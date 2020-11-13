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

package metrictree

import (
	"errors"
	"fmt"
)

func MatchNsToFilter(ns string, filter string) (bool, error) {
	parsedNs, err := ParseNamespace(ns, false)
	if err != nil {
		return false, err
	}
	if !parsedNs.IsUsableForAddition(true, false) {
		return false, fmt.Errorf("invalid format of the namespace: %v", ns)
	}

	// handle special case when filter have 1 element
	el, _, err := SplitNamespace(filter)
	if err != nil {
		return false, fmt.Errorf("invalid filter: %w", err)
	}

	var parsedFilter *Namespace
	if len(el) == 2 { // el[0] is empty
		if el[1] == "" {
			return true, nil
		}
		ne, errPNE := parseNamespaceElement(el[1], false)
		if errPNE != nil {
			return false, fmt.Errorf("invalid filter: %w", errPNE)
		}
		if _, ok := ne.(*staticSpecificElement); !ok {
			return false, fmt.Errorf("invalid filter: %v", filter)
		}

		parsedFilter = &Namespace{elements: []namespaceElement{ne}}
	} else {
		parsedFilter, err = ParseNamespace(filter, true)
		if err != nil {
			return false, err
		}

		if !parsedFilter.IsUsableForFiltering(true) {
			return false, fmt.Errorf("invalid format of the filter: %v", filter)
		}
	}

	if len(parsedFilter.elements) > len(parsedNs.elements) {
		return false, errors.New("namespace should have at least the same number of elements as filter")
	}

	for i, parsedFilterEl := range parsedFilter.elements {
		if !parsedFilterEl.Match(parsedNs.elements[i].String()) {
			return false, nil
		}
	}

	return true, nil
}
