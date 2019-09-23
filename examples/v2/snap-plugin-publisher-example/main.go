package main

import (
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
	return nil
}

func main() {
	runner.StartPublisher(&myPublisher{}, pluginName, pluginVersion)
}
