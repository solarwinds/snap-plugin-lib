package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/07-proxy/collector"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/07-proxy/collector/proxy"
)

/*****************************************************************************/

const pluginName = "system-collector"
const pluginVersion = "1.0.0"

func main() {
	runner.StartCollector(collector.New(proxy.New()), pluginName, pluginVersion)
}
