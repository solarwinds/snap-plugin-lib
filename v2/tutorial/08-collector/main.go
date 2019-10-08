package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	collector2 "github.com/librato/snap-plugin-lib-go/v2/tutorial/08-collector/collector"
	proxy2 "github.com/librato/snap-plugin-lib-go/v2/tutorial/08-collector/collector/proxy"
)

/*****************************************************************************/

const pluginName = "system-collector"
const pluginVersion = "1.0.0"

func main() {
	runner.StartCollector(collector2.New(proxy2.New()), pluginName, pluginVersion)
}
