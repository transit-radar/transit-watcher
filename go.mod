module github.com/catouberos/transit-watcher

go 1.24.3

toolchain go1.24.5

require github.com/google/uuid v1.6.0

require (
	buf.build/gen/go/catou/transit-radar/connectrpc/go v1.19.1-20251017072010-ae2a9f9d5b9c.2
	buf.build/gen/go/catou/transit-radar/protocolbuffers/go v1.36.10-20251017072010-ae2a9f9d5b9c.1
	connectrpc.com/connect v1.19.1
	github.com/cenkalti/backoff/v5 v5.0.3
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.14.0
	go.opentelemetry.io/otel/log v0.14.0
	go.opentelemetry.io/otel/sdk/log v0.14.0
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
