package metrictree

import (
	"fmt"
	"regexp"
)

type namespaceElement interface {
	Match(string) bool
	String() string
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

func newStaticRegexpElement(selector string) *staticRegexpElement {
	r, err := regexp.Compile(selector)
	if err != nil {
		// todo: log error
		return nil
	}
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
	return true // todo: calculate
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

func (dse *dynamicSpecificElement) Match(string) bool {
	return true // todo: calculate
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

func newDynamicRegexpElement(group, regText string) *dynamicRegexpElement {
	r, err := regexp.Compile(regText)
	if err != nil {
		// todo: log error
		return nil
	}

	return &dynamicRegexpElement{
		group:  group,
		regexp: r,
	}
}

func (dre *dynamicRegexpElement) Match(string) bool {
	return true // implement
}

func (dre *dynamicRegexpElement) String() string {
	return fmt.Sprintf("[%s={%s}]", dre.group, dre.regexp.String())
}
