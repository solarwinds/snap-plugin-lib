/*
 Copyright 2016 Intel Corporation

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

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

import "fmt"

var (
	// ErrEmptyKey is returned when a Rule with an empty key is created
	ErrEmptyKey = fmt.Errorf("Key cannot be Empty")

	// ErrConfigNotFound is returned when a config doesn't exist in the config map
	ErrConfigNotFound = fmt.Errorf("config item not found")

	// ErrNotA<type> is returned when the found config item doesn't have the expected type
	ErrNotAString = fmt.Errorf("config item is not a string")
	ErrNotAnInt   = fmt.Errorf("config item is not an int64")
	ErrNotABool   = fmt.Errorf("config item is not a boolean")
	ErrNotAFloat  = fmt.Errorf("config item is not a float64")
)
