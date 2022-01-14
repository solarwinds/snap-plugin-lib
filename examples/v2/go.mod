module github.com/solarwinds/snap-plugin-lib/examples/v2

go 1.13

require (
	github.com/josephspurrier/goversioninfo v1.4.0
	github.com/securego/gosec/v2 v2.9.5
	github.com/sirupsen/logrus v1.8.1
	github.com/solarwinds/snap-plugin-lib/v2 v2.0.4
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/tools v0.1.8
	honnef.co/go/tools v0.2.2
)

replace github.com/solarwinds/snap-plugin-lib/v2 => ../../v2
