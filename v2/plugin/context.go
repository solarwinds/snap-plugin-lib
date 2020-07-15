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

package plugin

import (
	"github.com/sirupsen/logrus"
)

// CollectContext provides state and configuration API to be used by custom code.
type Context interface {
	// Returns configuration value by providing path (representing its position in JSON tree)
	Config(key string) (string, bool)

	// Returns list of allowed configuration paths
	ConfigKeys() []string

	// Return raw configuration (JSON string)
	RawConfig() []byte

	// Store any object using key to have access from different Collect requests
	Store(key string, value interface{})

	// Load any object using key from different Collect requests (returns an interface{} which need to be casted to concrete type)
	Load(key string) (interface{}, bool)

	// Load any object using key from different Collect requests (passing it to provided reference).
	// Will throw error when dest type doesn't match to type of stored value or object with a given key wasn't found.
	LoadTo(key string, dest interface{}) error

	// Add warning information to current collect / process operation.
	AddWarning(msg string)

	// Check if task is completed
	IsDone() bool

	// Check if task is completed (via listening on a channel)
	Done() <-chan struct{}

	// Reference to logger object
	Logger() logrus.FieldLogger
}
