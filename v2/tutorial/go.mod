module github.com/solarwinds/snap-plugin-lib/v2/tutorial

go 1.13

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/josephspurrier/goversioninfo v1.4.0
	github.com/securego/gosec/v2 v2.9.5
	github.com/shirou/gopsutil v3.21.2+incompatible
	github.com/solarwinds/snap-plugin-lib/v2 v2.0.4
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/tools v0.1.8
	honnef.co/go/tools v0.2.2
)

replace github.com/solarwinds/snap-plugin-lib/v2 => ./..
