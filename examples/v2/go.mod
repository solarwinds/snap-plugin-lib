module github.com/solarwinds/snap-plugin-lib/examples/v2

go 1.13

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/josephspurrier/goversioninfo v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.0
	github.com/solarwinds/snap-plugin-lib/v2 v2.0.4
	github.com/urfave/cli v1.22.5 // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace github.com/solarwinds/snap-plugin-lib/v2 => ../../v2
