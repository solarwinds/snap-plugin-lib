module github.com/solarwinds/snap-plugin-lib/examples/v2

go 1.13

require (
	github.com/sirupsen/logrus v1.8.0
	github.com/solarwinds/snap-plugin-lib/v2 v2.0.4
	google.golang.org/grpc v1.36.0 // indirect
)

replace github.com/solarwinds/snap-plugin-lib/v2 => ../../v2
