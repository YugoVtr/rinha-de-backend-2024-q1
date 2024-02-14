FROM golang:1.22-alpine as builder

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/bin/app /
CMD ["/app"]
EXPOSE 8080
