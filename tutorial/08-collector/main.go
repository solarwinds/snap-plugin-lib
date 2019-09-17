package main

import (
	"github.com/librato/snap-plugin-lib-go/tutorial/08-collector/collector"
	"github.com/librato/snap-plugin-lib-go/tutorial/08-collector/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

/*****************************************************************************/

const pluginName = "system-collector"
const pluginVersion = "1.0.0"

func main() {
	runner.StartCollector(collector.New(proxy.New()), pluginName, pluginVersion)
}
