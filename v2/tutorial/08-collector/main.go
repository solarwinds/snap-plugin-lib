package main

import (
	"context"

	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/08-collector/collector"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/08-collector/collector/proxy"
)

/*****************************************************************************/

const pluginName = "system-collector"
const pluginVersion = "1.0.0"

func main() {
	runner.StartCollector(context.Background(), collector.New(proxy.New()), pluginName, pluginVersion)
}
