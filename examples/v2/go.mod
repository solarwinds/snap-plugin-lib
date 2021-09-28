module github.com/solarwinds/snap-plugin-lib/examples/v2

go 1.13

require (
	github.com/securego/gosec/v2 v2.8.1
	github.com/sirupsen/logrus v1.8.0
	github.com/solarwinds/snap-plugin-lib/v2 v2.0.4
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/tools v0.1.3
	honnef.co/go/tools v0.0.1-2020.1.4
)

replace github.com/solarwinds/snap-plugin-lib/v2 => ../../v2
