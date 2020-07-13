/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

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
