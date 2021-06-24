// +build tools

package tools

import (
	_ "github.com/a8m/envsubst/cmd/envsubst"
	_ "github.com/josephspurrier/goversioninfo/cmd/goversioninfo"
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
