package tracer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

func newResource(ctx context.Context, serviceNameUsedDisplayTracesBackends string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceNameUsedDisplayTracesBackends),
		),
	)
}

func connectsOpenTelemetryCollector(gRPCConnection string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		gRPCConnection,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func setupTraceExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

// Register the trace exporter with a TracerProvider, using a batch
// span processor to aggregate spans before export.
func registerTraceExporterUsingBatchSpanProcessor(exporter *otlptrace.Exporter, resource *resource.Resource) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider
}

func InitProvider(collectorURL, serviceName string) (func(context.Context) error, error) {
	slog.Info("tracer provider started", "collector URL", collectorURL)
	ctx := context.Background()
	res, err := newResource(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	conn, err := connectsOpenTelemetryCollector(collectorURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	traceExporter, err := setupTraceExporter(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	tracerProvider := registerTraceExporterUsingBatchSpanProcessor(traceExporter, res)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

type Trace struct {
	trace.Tracer
}

func New(name string) Trace {
	return Trace{otel.Tracer(name)}
}

func Middleware(service string) mux.MiddlewareFunc {
	return otelmux.Middleware(service)
}

func SetSpanError(err error) (codes.Code, string) {
	if err == nil {
		return codes.Ok, ""
	}
	return codes.Error, err.Error()
}

type SpanOpts map[string]any

func GetSpanOpts(opts []SpanOpts) []trace.SpanStartOption {
	spanOpts := make([]trace.SpanStartOption, 0, len(opts))
	for _, opt := range opts {
		for key, opt := range opt {
			spanOpts = append(spanOpts,
				trace.WithAttributes(
					attribute.String(key, fmt.Sprint(opt)),
				),
			)
		}
	}
	return spanOpts
}
