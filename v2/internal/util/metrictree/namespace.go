package metrictree

import (
	"fmt"
	"regexp"
	"strings"
)

type namespace struct {
	elements []namespaceElement
}

type namespaceElement interface {
	Match(string) bool
	String() string
}

/*****************************************************************************/

// Check is namespace selector can be used for metric definition
// First and last element should be static names, middle elements can be group (ie. [group])
func (ns *namespace) isUsableForDefinition() bool {
	if len(ns.elements) < 2 {
		return false
	}

	for _, nsElem := range []namespaceElement{ns.elements[0], ns.elements[len(ns.elements)-1]} {
		switch nsElem.(type) {
		case *staticSpecificElement: // ok
		default:
			return false
		}
	}

	for _, nsElem := range ns.elements[1 : len(ns.elements)-1] {
		switch nsElem.(type) {
		case *staticSpecificElement: // ok
		case *dynamicAnyElement: // ok
		default:
			return false
		}
	}

	return true
}

// Check is namespace selector can be used for metric addition ie. in ctx.AddMetric
// First and last element should be static names, middle elements can be group with defined value (ie. [group=id])
func (ns *namespace) isUsableForAddition() bool {
	if len(ns.elements) < 2 {
		return false
	}

	for _, nsElem := range []namespaceElement{ns.elements[0], ns.elements[len(ns.elements)-1]} {
		switch nsElem.(type) {
		case *staticSpecificElement: // ok
		default:
			return false
		}
	}

	for _, nsElem := range ns.elements[1 : len(ns.elements)-1] {
		switch nsElem.(type) {
		case *staticSpecificElement: // ok
		case *dynamicSpecificElement: // ok
		default:
			return false
		}
	}

	return true
}

func (ns *namespace) isUsableForSelection() bool {
	return true // todo: https://swicloud.atlassian.net/browse/AO-12232
}

/*****************************************************************************/

// Representing 2nd element of: /plugin/*/metric1
type staticAnyElement struct {
}

func newStaticAnyElement() *staticAnyElement {
	return &staticAnyElement{}
}

func (*staticAnyElement) Match(string) bool {
	return true
}

func (*staticAnyElement) String() string {
	return "*"
}

/*****************************************************************************/

// Representing 2nd element of: /plugin/group1/metric1
type staticSpecificElement struct {
	name string
}

func newStaticSpecificElement(selector string) *staticSpecificElement {
	return &staticSpecificElement{
		name: selector,
	}
}

func (sse *staticSpecificElement) Match(s string) bool {
	return sse.name == s
}

func (sse *staticSpecificElement) String() string {
	return sse.name
}

/*****************************************************************************/

// Representing 2nd element of: /plugin/{"group.*"}/metric1
type staticRegexpElement struct {
	regExp *regexp.Regexp
}

func newStaticRegexpElement(r *regexp.Regexp) *staticRegexpElement {
	return &staticRegexpElement{
		regExp: r,
	}
}

func (sre *staticRegexpElement) Match(s string) bool {
	return sre.regExp.MatchString(s)
}

func (sre *staticRegexpElement) String() string {
	return "{" + sre.regExp.String() + "}" // todo: constants
}

/*****************************************************************************/

// Representing 2nd element of: /plugin/[group1]/metric1
type dynamicAnyElement struct {
	group string
}

func newDynamicAnyElement(group string) *dynamicAnyElement {
	return &dynamicAnyElement{
		group: group,
	}
}

func (dae *dynamicAnyElement) Match(s string) bool {
	if isSurroundedWith(s, dynamicElementBeginIndicator, dynamicElementEndIndicator) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if len(dynElem) >= 3 && eqIndex > 0 && eqIndex < len(dynElem)-1 { // todo: remove duplication
				groupName := dynElem[0:eqIndex]
				return groupName == dae.group
			}
		}
	}

	return true
}

func (dae *dynamicAnyElement) String() string {
	return fmt.Sprintf("[%s]", dae.group)
}

/*****************************************************************************/

// Representing 2nd element of: /plugin/[group=id1]/metric1
type dynamicSpecificElement struct {
	group string
	value string
}

func (dse *dynamicSpecificElement) Match(s string) bool {
	if isSurroundedWith(s, dynamicElementBeginIndicator, dynamicElementEndIndicator) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if len(dynElem) >= 3 && eqIndex > 0 && eqIndex < len(dynElem)-1 { // todo: remove duplication
				groupName := dynElem[0:eqIndex]
				groupValue := dynElem[eqIndex+1:]

				return dse.group == groupName && dse.value == groupValue
			}
		}
	} else {
		if dse.value == s {
			return true
		}
	}

	return false
}

func (dse *dynamicSpecificElement) String() string {
	return fmt.Sprintf("[%s=%s]", dse.group, dse.value)
}

func newDynamicSpecificElement(group, value string) *dynamicSpecificElement {
	return &dynamicSpecificElement{
		group: group,
		value: value,
	}
}

/*****************************************************************************/

// Representing 2nd element of: /plugin/[group={id.*}/metric1
type dynamicRegexpElement struct {
	group  string
	regexp *regexp.Regexp
}

func newDynamicRegexpElement(group string, r *regexp.Regexp) *dynamicRegexpElement {
	return &dynamicRegexpElement{
		group:  group,
		regexp: r,
	}
}

func (dre *dynamicRegexpElement) Match(s string) bool {
	if isSurroundedWith(s, dynamicElementBeginIndicator, dynamicElementEndIndicator) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if len(dynElem) >= 3 && eqIndex > 0 && eqIndex < len(dynElem)-1 { // todo: remove duplication
				groupName := dynElem[0:eqIndex]
				groupValue := dynElem[eqIndex+1:]

				return dre.group == groupName && dre.regexp.MatchString(groupValue)
			}
		}
	} else {
		if dre.regexp.MatchString(s) {
			return true
		}
	}

	return false
}

func (dre *dynamicRegexpElement) String() string {
	return fmt.Sprintf("[%s={%s}]", dre.group, dre.regexp.String())
}
