package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	collector2 "github.com/librato/snap-plugin-lib-go/v2/tutorial/07-proxy/collector"
	proxy2 "github.com/librato/snap-plugin-lib-go/v2/tutorial/07-proxy/collector/proxy"
)

/*****************************************************************************/

const pluginName = "system-collector"
const pluginVersion = "1.0.0"

func main() {
	runner.StartCollector(collector2.New(proxy2.New()), pluginName, pluginVersion)
}
