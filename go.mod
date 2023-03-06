module github.com/solarwinds/snap-plugin-lib

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/google/go-licenses v0.0.0-00010101000000-000000000000
	github.com/josephspurrier/goversioninfo v1.4.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/securego/gosec/v2 v2.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/cobra v1.3.0 // indirect
	github.com/urfave/cli v1.22.5
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/tools v0.1.12
	google.golang.org/genproto v0.0.0-20220112215332-a9c7c0acf9f2 // indirect
	google.golang.org/grpc v1.43.0
	honnef.co/go/tools v0.3.2
)

replace (
	github.com/google/go-licenses => github.com/google/go-licenses v0.0.0-20211216151313-baa4dc16ea37

	// Freeze as in the next commit there was //go:embed added (supported since go 1.16)
	github.com/google/licenseclassifier => github.com/google/licenseclassifier v0.0.0-20210325184830-bb04aff29e72
	github.com/securego/gosec/v2 => github.com/securego/gosec/v2 v2.9.1
)
