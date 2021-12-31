/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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

package runner

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func printVersion(name string, version string) {
	fmt.Printf("%s version %s\n", name, version)

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Printf("\tbuilt with:\n")

		for _, dep := range buildInfo.Deps {
			if strings.Contains(dep.Path, "snap-") {
				path := strings.SplitAfter(dep.Path, "snap-")[1]
				path = strings.Split(path, "/")[0]
				fmt.Printf("\t%v (%v)\n", path, dep.Version)
			}
		}
	}

}
