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

package types

import (
	"fmt"
	"strings"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

type Namespace []NamespaceElement

func (ns Namespace) HasElement(el string) bool {
	for _, nsElem := range ns {
		if el == nsElem.String() {
			return true
		}
	}

	return false
}

func (ns Namespace) HasElementOn(el string, pos int) bool {
	if pos < len(ns) && pos >= 0 {
		if el == ns[pos].String() {
			return true
		}
	}

	return false
}

func (ns Namespace) String() string {
	var sb strings.Builder

	sb.WriteString(metricSeparator)

	for i, nsElem := range ns {
		sb.WriteString(nsElem.String())

		if i != len(ns)-1 {
			sb.WriteString(metricSeparator)
		}
	}

	return sb.String()
}

func (ns Namespace) Len() int {
	return len(ns)
}

func (ns Namespace) At(pos int) plugin.NamespaceElement {
	return &ns[pos]
}

///////////////////////////////////////////////////////////////////////////////

type NamespaceElement struct {
	Name_        string
	Value_       string
	Description_ string
}

func (ns *NamespaceElement) Name() string {
	return ns.Name_
}

func (ns *NamespaceElement) Value() string {
	return ns.Value_
}

func (ns *NamespaceElement) Description() string {
	return ns.Description_
}

func (ns *NamespaceElement) IsDynamic() bool {
	return ns.Name_ != ""
}

func (ns *NamespaceElement) String() string {
	if ns.Name_ == "" {
		return ns.Value_
	}

	return fmt.Sprintf("[%s=%s]", ns.Name_, ns.Value_)
}
