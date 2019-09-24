/*
This file contains definition of namespace and namespace elements. amespace element can consist of:
- static names
- groups - elements which can hold additional values (and are converted to tags by AO)
- regular expressions - used to restrict metrics sent to AO
- special selectors - * and **
*/

// Example of valid namespaces:
//    /plugin/group1/metric1             - general metric (or general filter)
//    /plugin/[group1]/metric1           - metric containing group
//    /plugin/[group1=id1]/metric1       - metric containing group with value
//    /plugin/{id.*}/metric1             - filter with regular expression matcher
//    /plugin/[group1={id.*}]/metric1    - filter with regular expression matcher for groups
//    /plugin/group1/*/metric1           - filter matching anything to single namespace element
//    /plugin/group1/**                  - filter matching anything to namespace prefix (** match to 1 or more namespace elements)

/*
Namespaces (very often referred as "selectors") are used in three different contexts:
    - when defining metrics (PluginDefinition) by plugin
    - when defining filters (requested metrics) (task*.yaml)
    - when adding concrete metrics during collection (CollectMetrics)

Not all forms are valid for different context. Ie.
    /plugin/[group]/metric
is a valid form of metric definition and filter, but can't be used to add metric since it's not concrete

On the other hand:
    /plugin/[group=id1]/metric
is a valid form of filter and can be used for addition, but it's not acceptable as a definition (group can't be concrete)

Look at IsUsableFor*() and tests to understand all possible cases.
*/

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
	// used to validate "concrete metric" against "filters" and "definitions" (should add metric to collect result?).
	// ie. "/plugin/[node=node123]/metric1" is valid against filter "/plugin/**" and definition "/plugin/[node]/metric1".
	Match(string) bool

	// check if "filters" match to "definition" (is given filter correct?).
	// Used to throw away filters which wouldn't have any impact on collection at early stage (loading task).
	// ie. Filters "/plugin/**", "/plugin/*/metric1", "/plugin/{node[1-3]{1,2}}" and compatible with definition "/plugin/[node]/metric1".
	Compatible(string) bool

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

// Check if namespace selector can be used for metric addition (ctx.AddMetric) or metric calculation reasonableness (ctx.ShouldProcess)
// First and last element should be static names, middle elements can be group with defined value (ie. [group=id])
//
// metricDefinitionPresent - In case plugin doesn't provide metric definition, added elements should be only static names.
// allowAnyMatch - When true, using '*' is allowed (ie. ctx.ShouldProcess("/plugin/group/*/*/metric1")
func (ns *Namespace) IsUsableForAddition(metricDefinitionPresent bool, allowAnyMatch bool) bool {
	if len(ns.elements) < minNamespaceLength {
		return false
	}

	if (allowAnyMatch && !ns.isFirstElementStatic()) || (!allowAnyMatch && !ns.isFirstAndLastElementStatic()) {
		return false
	}

	for _, nsElem := range ns.elements[1 : len(ns.elements)-1] {
		switch nsElem.(type) {
		case *staticAnyElement:
			if !allowAnyMatch {
				return false
			}
		case *staticSpecificElement: // ok
		case *dynamicSpecificElement:
			if !metricDefinitionPresent {
				return false
			}
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

func (ns *Namespace) isFirstElementStatic() bool {
	switch ns.elements[0].(type) {
	case *staticSpecificElement: // ok
	default:
		return false
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

func (*staticAnyElement) Match(s string) bool {
	return isValidIdentifier(s)
}

func (sae *staticAnyElement) Compatible(s string) bool {
	return false
}

func (*staticAnyElement) String() string {
	return staticAnyMatcher
}

func (*staticAnyElement) IsDynamic() bool { return false }
func (*staticAnyElement) HasRegexp() bool { return false }

/*****************************************************************************/

// Representing 3rd element of: /plugin/group1/**
type staticRecursiveAnyElement struct {
}

func newStaticRecursiveAnyElement() *staticRecursiveAnyElement {
	return &staticRecursiveAnyElement{}
}

func (*staticRecursiveAnyElement) Match(s string) bool {
	return isValidIdentifier(s)
}

func (*staticRecursiveAnyElement) Compatible(s string) bool {
	return false
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

func newStaticSpecificElement(name string) *staticSpecificElement {
	return &staticSpecificElement{
		name: name,
	}
}

func (sse *staticSpecificElement) Match(s string) bool {
	return sse.name == s
}

func (sse *staticSpecificElement) Compatible(s string) bool {
	if containsGroup(s) {
		return false
	}

	if containsRegexp(s) || s == staticAnyMatcher || s == staticRecursiveAnyMatcher {
		return true
	}

	if s == sse.name {
		return true
	}

	return false
}

func (sse *staticSpecificElement) String() string {
	return sse.name
}

func (*staticSpecificElement) IsDynamic() bool { return false }
func (*staticSpecificElement) HasRegexp() bool { return false }

/*****************************************************************************/

// Representing 2nd element of: /plugin/{group.*}/metric1
type staticRegexpElement struct {
	regExp *regexp.Regexp
}

func newStaticRegexpElement(r *regexp.Regexp) *staticRegexpElement {
	return &staticRegexpElement{
		regExp: r,
	}
}

func (sre *staticRegexpElement) Match(s string) bool {
	return !containsGroup(s) && isValidIdentifier(s) && sre.regExp.MatchString(s)
}

func (*staticRegexpElement) Compatible(s string) bool {
	return false
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
		eqIndex := strings.Index(dynElem, dynamicElementEqualIndicator)

		if eqIndex != -1 {
			groupName := dynElem[0:eqIndex]
			groupValue := dynElem[eqIndex+1:]

			return groupName == dae.group && isValidIdentifier(groupValue)
		}
	}

	return isValidIdentifier(s)
}

func (dae *dynamicAnyElement) Compatible(s string) bool {
	if containsGroup(s) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, dynamicElementEqualIndicator)

		groupName := dynElem
		if eqIndex != -1 {
			groupName = dynElem[0:eqIndex]
		}

		return groupName == dae.group
	}

	return true
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

func newDynamicSpecificElement(group, value string) *dynamicSpecificElement {
	return &dynamicSpecificElement{
		group: group,
		value: value,
	}
}

func (dse *dynamicSpecificElement) Match(s string) bool {
	if containsGroup(s) {
		dynElem := s[1 : len(s)-1]
		eqIndex := strings.Index(dynElem, dynamicElementEqualIndicator)

		if eqIndex != -1 {
			groupName := dynElem[0:eqIndex]
			groupValue := dynElem[eqIndex+1:]

			return dse.group == groupName && dse.value == groupValue
		}
	} else {
		if dse.value == s {
			return true
		}
	}

	return false
}

func (*dynamicSpecificElement) Compatible(s string) bool {
	return false
}

func (dse *dynamicSpecificElement) String() string {
	return fmt.Sprintf("[%s=%s]", dse.group, dse.value)
}

func (*dynamicSpecificElement) IsDynamic() bool { return false }
func (*dynamicSpecificElement) HasRegexp() bool { return true }

/*****************************************************************************/

// Representing 2nd element of: /plugin/[group={id.*}]/metric1
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
		eqIndex := strings.Index(dynElem, dynamicElementEqualIndicator)

		if eqIndex != -1 {
			groupName := dynElem[0:eqIndex]
			groupValue := dynElem[eqIndex+1:]

			return dre.group == groupName && isValidIdentifier(groupValue) && dre.regexp.MatchString(groupValue)
		}
	} else {
		if dre.regexp.MatchString(s) {
			return true
		}
	}

	return false
}

func (*dynamicRegexpElement) Compatible(s string) bool {
	return false
}

func (dre *dynamicRegexpElement) String() string {
	return fmt.Sprintf("[%s={%s}]", dre.group, dre.regexp.String())
}

func (*dynamicRegexpElement) IsDynamic() bool { return true }
func (*dynamicRegexpElement) HasRegexp() bool { return true }

/*****************************************************************************/

/*
Special case: representing 2nd element of: /plugin/group1/metric1 in filters when plugin provides definition. Ie.
We have metric defined:
	/plugin/[dyn1]/metric1
We can add filter using two methods:
	/plugin/id1/metric1
	/plugin/[dyn1=id]/metric1
*/
type staticSpecificAcceptingGroupElement struct {
	name string
}

func newStaticSpecificAcceptingGroupElement(name string) *staticSpecificAcceptingGroupElement {
	return &staticSpecificAcceptingGroupElement{
		name: name,
	}
}

func (sse *staticSpecificAcceptingGroupElement) Match(s string) bool {
	if sse.name == s {
		return true
	}

	parsedEl, err := parseNamespaceElement(s, false)
	if err != nil {
		return false
	}
	if gps, ok := parsedEl.(*dynamicSpecificElement); ok {
		return gps.value == sse.name
	}

	return false
}

func (*staticSpecificAcceptingGroupElement) Compatible(s string) bool {
	return false
}

func (sse *staticSpecificAcceptingGroupElement) String() string {
	return sse.name
}

func (*staticSpecificAcceptingGroupElement) IsDynamic() bool { return false }
func (*staticSpecificAcceptingGroupElement) HasRegexp() bool { return false }

/*****************************************************************************/

/*
Special case: representing 2nd element of: /plugin/{group1}/metric1 in filters when plugin provides definition. Ie.
We have metric defined:
	/plugin/[dyn1]/metric1
We can add filter using two methods:
	/plugin/{id1.*}/metric1
	/plugin/[dyn1={id.*}]/metric1
*/
type staticRegexpAcceptingGroupElement struct {
	regExp *regexp.Regexp
}

func newStaticRegexpAcceptingGroupElement(r *regexp.Regexp) *staticRegexpAcceptingGroupElement {
	return &staticRegexpAcceptingGroupElement{
		regExp: r,
	}
}

func (sre *staticRegexpAcceptingGroupElement) Match(s string) bool {
	parsedEl, err := parseNamespaceElement(s, false)
	if err != nil {
		return false
	}

	switch parsedEl.(type) {
	case *staticSpecificElement:
		return sre.regExp.MatchString(s)
	case *dynamicSpecificElement:
		return sre.regExp.MatchString(parsedEl.(*dynamicSpecificElement).value)
	}

	return false
}

func (*staticRegexpAcceptingGroupElement) Compatible(s string) bool {
	return false
}

func (sre *staticRegexpAcceptingGroupElement) String() string {
	return regexBeginIndicator + sre.regExp.String() + regexEndIndicator
}

func (*staticRegexpAcceptingGroupElement) IsDynamic() bool { return false }
func (*staticRegexpAcceptingGroupElement) HasRegexp() bool { return false }
