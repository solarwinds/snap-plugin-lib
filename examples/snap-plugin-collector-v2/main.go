package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin/runner"
	"github.com/librato/snap-plugin-lib-go/v2/plugin/types"
)

type myCollector struct {
}

func (*myCollector) Collect(ctx types.Context) error {
	panic("implement me")
}

func main() {
	runner.StartCollector(&myCollector{}, "my collector", "0.0.1")
}
