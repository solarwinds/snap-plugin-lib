module github.com/solarwinds/snap-plugin-lib/v2/tutorial

go 1.13

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/josephspurrier/goversioninfo v1.2.0 // indirect
	github.com/shirou/gopsutil v3.21.2+incompatible
	github.com/solarwinds/snap-plugin-lib/v2 v2.0.4
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	github.com/urfave/cli v1.22.5 // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace github.com/solarwinds/snap-plugin-lib/v2 => ./..
