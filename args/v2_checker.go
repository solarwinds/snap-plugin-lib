package args

import (
	"os"
	"strings"
)

// Check for presence of --plugin-api-v2 flag that changes the default operating mode to using plugin API v2
func UsePluginAPIv2() bool {
	for _, a := range os.Args {
		if strings.Contains(a, "-plugin-api-v2") {
			return true
		}
	}

	return false
}
