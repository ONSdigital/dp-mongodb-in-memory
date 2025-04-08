module github.com/ONSdigital/dp-mongodb-in-memory

go 1.24

require (
	github.com/ONSdigital/log.go/v2 v2.4.4
	github.com/smartystreets/goconvey v1.8.1
	github.com/spf13/afero v1.12.0
	go.mongodb.org/mongo-driver v1.17.2
	golang.org/x/crypto v0.32.0
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.263.0 // indirect
	github.com/ONSdigital/dp-net/v3 v3.0.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

retract (
	v1.4.0 // Now considered unsafe as the semantics of the Start() function implicitly changed
	v1.0.0 // Licensing error DO NOT USE
)
