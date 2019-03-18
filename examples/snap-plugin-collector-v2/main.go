package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

type myCollector struct {
}

func (*myCollector) Collect(ctx plugin.Context) error {
	panic("implement me")
}

func main() {
	runner.StartCollector(&myCollector{}, "my collector", "0.0.1")
}
