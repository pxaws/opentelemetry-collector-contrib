module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsprometheusremotewriteexporter

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.69
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.29.1-0.20210702174635-c64d1f096bda
	go.uber.org/zap v1.18.1
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)

replace go.opentelemetry.io/collector/model => go.opentelemetry.io/collector/model v0.0.0-20210702174635-c64d1f096bda
