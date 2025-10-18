module github.com/catouberos/transit-watcher

go 1.24.3

toolchain go1.24.5

require github.com/google/uuid v1.6.0

require (
	buf.build/gen/go/catou/transit-radar/connectrpc/go v1.19.1-20251017072010-ae2a9f9d5b9c.2
	buf.build/gen/go/catou/transit-radar/protocolbuffers/go v1.36.10-20251017072010-ae2a9f9d5b9c.1
	connectrpc.com/connect v1.19.1
	github.com/cenkalti/backoff/v5 v5.0.3
)

require google.golang.org/protobuf v1.36.10 // indirect
