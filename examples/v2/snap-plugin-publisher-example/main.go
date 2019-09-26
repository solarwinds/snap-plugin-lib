package main

import (
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

const (
	pluginName    = "example"
	pluginVersion = "0.0.1"
)

type myPublisher struct {
}

func (m myPublisher) Publish(ctx plugin.PublishContext) error {
	fmt.Printf("Number of metrics: %v\n", ctx.Count())

	for _, mt := range ctx.ListAllMetrics() {
		fmt.Printf(" - %s=%v [%v]\n", mt.Namespace(), mt.Value(), mt.Tags())
	}

	return nil
}

func main() {
	runner.StartPublisher(&myPublisher{}, pluginName, pluginVersion)
}
