//go:build small
// +build small

/*
 Copyright (c) 2022 SolarWinds Worldwide, LLC

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

package metrictree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	ns          string
	filter      string
	shouldMatch bool
	expectError bool
}

var testCases = []testCase{
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/kubernetes/[pod=io1453]/cpu",
		shouldMatch: true,
	},
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/kubernetes/[pod=io1453]",
		shouldMatch: true,
	},
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/kubernetes/[pod]",
		shouldMatch: true,
	},
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/kubernetes/[node]",
		shouldMatch: false,
	},
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/kubernetes/io1453",
		shouldMatch: true,
	},
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/kubernetes",
		shouldMatch: true,
	},
	{
		ns:          "/kubernetes/[pod=io1453]/cpu",
		filter:      "/",
		shouldMatch: true,
	},
	{
		ns:          "|kubernetes|[pod=io1453]|cpu",
		filter:      "|kubernetes|[pod={io[0-9]{4}}]|cpu",
		shouldMatch: true,
	},
	{
		ns:          "/kubernetes//cpu",
		filter:      "/kubernetes/[pod={io[0-9]{4}}]/cpu",
		expectError: true,
	},
	{
		ns:          "|log-files|[file=/tmp/sample.log]|string_line",
		filter:      "|log-files|[file=/tmp/sample.log]",
		shouldMatch: true,
	},
	{
		ns:          "|log-files|[file=/tmp/sample1.log]|string_line",
		filter:      "|log-files|[file=/tmp/sample.log]|string_line",
		shouldMatch: false,
	},
	{
		ns:          "/syslog/[ip_address=127.0.0.1]/[hostname=localhost]/string_line",
		filter:      "/syslog/*/[hostname=localhost]/string_line",
		shouldMatch: true,
	},
	{
		ns:          "/syslog/[ip_address=127.0.0.1]/[hostname=localhost]/string_line",
		filter:      "/syslog/[ip_address=127.0.0.1]/[hostname=localhost]/string_line",
		shouldMatch: true,
	},
	{
		ns:          "/syslog/[ip_address=127.0.1.1]/[hostname=localhost]/string_line",
		filter:      "/syslog/[ip_address=127.0.0.1]/[hostname=localhost]/string_line",
		shouldMatch: false,
	},
	{
		ns:          "/syslog/[ip_address=127.0.0.1]/[hostname=localhost]/string_line",
		filter:      "/syslog/**/[hostname=localhost]/string_line",
		shouldMatch: false,
		expectError: true,
	},
	{
		ns:          "/syslog/[ip_address=127.0.0.1]/[hostname=localhost]/string_line",
		filter:      "/syslog/[ip_address]/[hostname=localhost]/string_line",
		shouldMatch: true,
	},
	{
		ns:          "/syslog/[ip_address=127.0.0.1]/[hostname=otherhost]/string_line",
		filter:      "/syslog/*/[hostname=localhost]/string_line",
		shouldMatch: false,
	},
	{
		ns:          "/mock/metric/c",
		filter:      "/another-plugin",
		shouldMatch: false,
	},
	{
		ns:          "/mock/metric/c",
		filter:      "/",
		shouldMatch: true,
	},
	{
		ns:          "/mock/metric/c",
		filter:      "/*",
		shouldMatch: true,
	},
	{
		ns:          "/mock/metric/c",
		filter:      "/*/metric/d",
		shouldMatch: false,
	},
	{
		ns:          "/mock/metric/c",
		filter:      "/**",
		shouldMatch: true,
	},
	{
		ns:          "/mock/metric/c",
		filter:      "/mock",
		shouldMatch: true,
	},
}

func TestMatchNsToFilter(t *testing.T) {
	for i, te := range testCases {
		fmt.Printf("i=%#v\n", i)
		match, err := MatchNsToFilter(te.ns, te.filter)

		if te.shouldMatch {
			require.True(t, match, "Namespace should match to a given filter")
		} else {
			require.False(t, match, "Namespace shouldn't match to a given filter")
		}

		if te.expectError {
			require.Error(t, err, "Expected error")
		} else {
			require.NoError(t, err, "Expected no error")
		}
	}
}
