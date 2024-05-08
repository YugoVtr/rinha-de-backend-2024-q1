load('ext://uibutton', 'cmd_button', 'location')

local_resource(
  'webserver-api-build',
  cmd='make build',
  deps=['main.go', 'go.mod', "entity", "infra", "persistence", "server", "test"],
  labels=['api']
)

docker_build(
  'rinha-backend-image',
  '.',
  dockerfile='Dockerfile',
  only=['bin/out/rinha-de-backend-2024-q1'],
  live_update=[
    sync('./bin/out/rinha-de-backend-2024-q1', '/app/rinha-de-backend-2024-q1')
  ]
)

docker_compose('./compose.yml')
dc_resource('api01', labels=["api"])
dc_resource('api02', labels=["api"])

dc_resource('nginx', labels=["infra"])
dc_resource('db', labels=["infra"])
dc_resource('jaeger', labels=["infra"])
dc_resource('otel-collector', labels=["infra"])

cmd_button('reset',
  argv=['make', 'restart'],
  location=location.NAV,
  icon_name='refresh',
  text='Reset DB',
)

cmd_button('load test',
  argv=['make', 'load-test'],
  location=location.NAV,
  icon_name='rocket_launch',
  text='Load Test',
)
