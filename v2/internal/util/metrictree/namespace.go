package metrictree

import (
	"fmt"
	"regexp"
	"strings"
)

type Namespace struct {
	elements []namespaceElement
}

type namespaceElement interface {
	Match(string) bool
	String() string

	IsDynamic() bool
	HasRegexp() bool
}

const minNamespaceLength = 2

/*****************************************************************************/

// Check is namespace selector can be used for metric definition
// First and last element should be static names, middle elements can be group (ie. [group])
func (ns *Namespace) IsUsableForDefinition() bool {
	if len(ns.elements) < minNamespaceLength {
		return false
	}

	if !ns.isFirstAndLastElementStatic() {
		return false
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

// Check if namespace selector can be used for metric addition ie. in ctx.AddMetric
// First and last element should be static names, middle elements can be group with defined value (ie. [group=id])
func (ns *Namespace) IsUsableForAddition() bool {
	if len(ns.elements) < minNamespaceLength {
		return false
	}

	if !ns.isFirstAndLastElementStatic() {
		return false
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

// Check if namespace selector can be used for metric filters
// !! Note: If metric definition is not provided in plugin, matcher with dynamic element can't be used in filter (to avoid ambiguity)
func (ns *Namespace) IsUsableForFiltering(metricDefinitionPresent bool) bool {
	if len(ns.elements) < minNamespaceLength {
		return false
	}

	switch ns.elements[0].(type) {
	case *staticSpecificElement:
	case *staticSpecificAcceptingGroupElement:
	default:
		return false
	}

	if !metricDefinitionPresent {
		for _, nsElem := range ns.elements[1:len(ns.elements)] {
			if nsElem.IsDynamic() == true {
				return false
			}
		}
	}

	return true
}

func (ns *Namespace) isFirstAndLastElementStatic() bool {
	for _, nsElem := range []namespaceElement{ns.elements[0], ns.elements[len(ns.elements)-1]} {
		switch nsElem.(type) {
		case *staticSpecificElement: // ok
		default:
			return false
		}
	}

	return true
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
	return string(staticAnyMatcher)
}

func (*staticAnyElement) IsDynamic() bool { return false }
func (*staticAnyElement) HasRegexp() bool { return false }

/*****************************************************************************/

// Representing 2nd element of: /plugin/group1/**
type staticRecursiveAnyElement struct {
}

func newStaticRecursiveAnyElement() *staticRecursiveAnyElement {
	return &staticRecursiveAnyElement{}
}

func (*staticRecursiveAnyElement) Match(string) bool {
	return true
}

func (*staticRecursiveAnyElement) String() string {
	return staticRecursiveAnyMatcher
}

func (*staticRecursiveAnyElement) IsDynamic() bool { return false }
func (*staticRecursiveAnyElement) HasRegexp() bool { return false }

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

func (*staticSpecificElement) IsDynamic() bool { return false }
func (*staticSpecificElement) HasRegexp() bool { return false }

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
	return !containsGroup(s) && sre.regExp.MatchString(s)
}

func (sre *staticRegexpElement) String() string {
	return regexBeginIndicator + sre.regExp.String() + regexEndIndicator
}

func (*staticRegexpElement) IsDynamic() bool { return false }
func (*staticRegexpElement) HasRegexp() bool { return true }

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
	if containsGroup(s) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if isIndexInTheMiddle(eqIndex, s) {
				groupName := dynElem[0:eqIndex]
				groupValue := dynElem[eqIndex+1:]

				return groupName == dae.group && isValidIdentifier(groupValue)
			}
		}
	}

	return isValidIdentifier(s)
}

func (dae *dynamicAnyElement) String() string {
	return fmt.Sprintf("[%s]", dae.group)
}

func (*dynamicAnyElement) IsDynamic() bool { return true }
func (*dynamicAnyElement) HasRegexp() bool { return false }

/*****************************************************************************/

// Representing 2nd element of: /plugin/[group=id1]/metric1
type dynamicSpecificElement struct {
	group string
	value string
}

func (dse *dynamicSpecificElement) Match(s string) bool {
	if containsGroup(s) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if isIndexInTheMiddle(eqIndex, s) {
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

func (*dynamicSpecificElement) IsDynamic() bool { return false }
func (*dynamicSpecificElement) HasRegexp() bool { return true }

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
	if containsGroup(s) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, string(dynamicElementEqualIndicator))

		if eqIndex != -1 {
			if isIndexInTheMiddle(eqIndex, s) {
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

func (*dynamicRegexpElement) IsDynamic() bool { return true }
func (*dynamicRegexpElement) HasRegexp() bool { return true }

/*****************************************************************************/

type staticSpecificAcceptingGroupElement struct {
	name string
}

func newStaticSpecificAcceptingGroupElement(selector string) *staticSpecificAcceptingGroupElement {
	return &staticSpecificAcceptingGroupElement{
		name: selector,
	}
}

func (sse *staticSpecificAcceptingGroupElement) Match(s string) bool {
	if sse.name == s {
		return true
	}

	ps, err := parseNamespaceElement(s, false)
	if err != nil {
		return false
	}
	if gps, ok := ps.(*dynamicSpecificElement); ok {
		return gps.value == sse.name
	}

	return false
}

func (sse *staticSpecificAcceptingGroupElement) String() string {
	return sse.name
}

func (*staticSpecificAcceptingGroupElement) IsDynamic() bool { return false }
func (*staticSpecificAcceptingGroupElement) HasRegexp() bool { return false }
