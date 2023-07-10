module github.com/ONSdigital/dp-mongodb-in-memory

go 1.19

// to fix: [CVE-2023-32731] CWE-Other
replace google.golang.org/grpc => google.golang.org/grpc v1.55.0

require (
	github.com/ONSdigital/log.go/v2 v2.4.1
	github.com/smartystreets/goconvey v1.8.0
	github.com/spf13/afero v1.9.5
	go.mongodb.org/mongo-driver v1.11.4
	golang.org/x/crypto v0.10.0
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.252.0 // indirect
	github.com/ONSdigital/dp-net/v2 v2.9.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/smartystreets/assertions v1.13.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
)

retract (
	v1.4.0 // Now considered unsafe as the semantics of the Start() function implicitly changed
	v1.0.0 // Licensing error DO NOT USE
)
