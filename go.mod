module github.com/catouberos/transit-watcher

go 1.24.3

toolchain go1.24.5

require github.com/google/uuid v1.6.0

require (
	buf.build/gen/go/catou/transit-radar/connectrpc/go v1.19.1-20251224091833-811a6157b8b0.2
	buf.build/gen/go/catou/transit-radar/protocolbuffers/go v1.36.11-20251224091833-811a6157b8b0.1
	connectrpc.com/connect v1.19.1
	github.com/cenkalti/backoff/v5 v5.0.3
	github.com/stretchr/testify v1.11.1
	github.com/twmb/franz-go v1.20.6
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.14.0
	go.opentelemetry.io/otel/log v0.14.0
	go.opentelemetry.io/otel/sdk/log v0.14.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto v0.0.0-20251222181119-0a764e51fe1b // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/spf13/viper v1.21.0
	github.com/twmb/franz-go/pkg/kmsg v1.12.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	google.golang.org/protobuf v1.36.11
)
