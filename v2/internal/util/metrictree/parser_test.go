// +build small

package metrictree

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*****************************************************************************/

type parseNamespaceValidScenario struct {
	namespace                         string
	usableForDefinition               bool
	usableForAdditionWhenDefinition   bool
	usableForAdditionWhenNoDefinition bool
	usableForSelection                bool
}

var parseNamespaceValidScenarios = []parseNamespaceValidScenario{
	{
		namespace:                         "/plugin/group1/metric",
		usableForDefinition:               true,
		usableForAdditionWhenDefinition:   true,
		usableForAdditionWhenNoDefinition: true,
	},
	{
		namespace:                         "/plugin/[group2]/metric",
		usableForDefinition:               true,
		usableForAdditionWhenDefinition:   false,
		usableForAdditionWhenNoDefinition: false,
	},
	{ // 2
		namespace:                         "/plugin/[group2=id]/metric",
		usableForDefinition:               false,
		usableForAdditionWhenDefinition:   true,
		usableForAdditionWhenNoDefinition: false,
	},
	{ // 3
		namespace:                         "/plugin/[group2={id.*}]/metric",
		usableForDefinition:               false,
		usableForAdditionWhenDefinition:   false,
		usableForAdditionWhenNoDefinition: false,
	},
	{ // 4
		namespace:                         "/plugin/{id.*}/metric",
		usableForDefinition:               false,
		usableForAdditionWhenDefinition:   false,
		usableForAdditionWhenNoDefinition: false,
	},
	{ // 5
		namespace:                         "/plugin/*/metric",
		usableForDefinition:               false,
		usableForAdditionWhenDefinition:   false,
		usableForAdditionWhenNoDefinition: false,
	},
	{ // 6
		namespace:                         "/plugin/metric/**",
		usableForDefinition:               false,
		usableForAdditionWhenDefinition:   false,
		usableForAdditionWhenNoDefinition: false,
	},
	{ // 7
		namespace:                         "/plugin/metric",
		usableForDefinition:               true,
		usableForAdditionWhenDefinition:   true,
		usableForAdditionWhenNoDefinition: true,
	},
	{ // 8
		namespace:                         "/plugin/**",
		usableForDefinition:               false,
		usableForAdditionWhenDefinition:   false,
		usableForAdditionWhenNoDefinition: false,
	},
}

func TestParseNamespace_ValidScenarios(t *testing.T) {
	Convey("Validate ParseNamespace - valid scenarios", t, func() {
		for i, tc := range parseNamespaceValidScenarios {
			Convey(fmt.Sprintf("Scenario %d", i), func() {
				// Act
				ns, err := ParseNamespace(tc.namespace, false)

				// Assert
				So(ns, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(ns.IsUsableForDefinition(), ShouldEqual, tc.usableForDefinition)
				So(ns.IsUsableForAddition(false), ShouldEqual, tc.usableForAdditionWhenNoDefinition)
				So(ns.IsUsableForAddition(true), ShouldEqual, tc.usableForAdditionWhenDefinition)
				So(ns.IsUsableForFiltering(true), ShouldBeTrue)
			})
		}
	})
}

/*****************************************************************************/

func TestParseNamespace_InvalidScenarios(t *testing.T) {
	testCases := []string{
		"/",
		"el",
		"el1/",
		"/el1/el2//m3",
		"/el1/el_#/m4",
		"el/el2/el3/el4",
		"/el/el2/el3/el4/",
		"/el/el2/**/m4",
		"/el/el2/**/",
		"/el/el2/ gr/m1",
		"/el/el2/gr /m1",
		"/el/el2/gr/ m1",
		"/el/el2/gr/ m1 ",
	}

	Convey("Validate ParseNamespace - negative scenarios", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("Scenario %d (%s)", i, tc), func() {
				// Act
				_, err := ParseNamespace(tc, false)

				// Assert
				So(err, ShouldBeError)
			})
		}
	})
}

/*****************************************************************************/

type parseNamespaceElementValidScenario struct {
	namespaceElement string
	comparableType   namespaceElement
	shouldMatch      []string
	shouldNotMatch   []string
	isFilter         bool
}

var parseNamespaceElementValidScenarios = []parseNamespaceElementValidScenario{
	{
		namespaceElement: "[group]",
		comparableType:   &dynamicAnyElement{},
		shouldMatch:      []string{"[group=id1]", "[group=id3]", "id3", "group"},
		shouldNotMatch:   []string{"[grp=id1]", "[group]", "[id1=group]", "*", "**", "{group}", ""},
		isFilter:         false,
	},
	{
		namespaceElement: "[group=id1]",
		comparableType:   &dynamicSpecificElement{},
		shouldMatch:      []string{"id1", "[group=id1]"},
		shouldNotMatch:   []string{"id2", "[group=id2]", "[grp=id1]", "[group]", "*", "**", "{group}", ""},
		isFilter:         false,
	},
	{
		namespaceElement: "[group={id.*}]",
		comparableType:   &dynamicRegexpElement{},
		shouldMatch:      []string{"[group=id1]", "[group=id3]", "id1", "id3"},
		shouldNotMatch:   []string{"[group=i1]", "[grp=id1]", "i1", "*", "**", "{group}", "[group={id1}]"},
		isFilter:         false,
	},
	{
		namespaceElement: "{}", // valid
		comparableType:   &staticRegexpElement{},
		shouldMatch:      []string{"id"},
		shouldNotMatch:   []string{"", "*", "**", "{group}", "[group=id]"},
		isFilter:         false,
	},
	{
		namespaceElement: "{mem.*[1-3]{1,}}",
		comparableType:   &staticRegexpElement{},
		shouldMatch:      []string{"memory3", "mem1", "memo2"},
		shouldNotMatch:   []string{"memo4", "memory0", "group", "[grp=memory3]", "", "*", "**", "{memo2}", "[group={mem1}]"},
		isFilter:         false,
	},
	{
		namespaceElement: "*", // valid
		comparableType:   &staticAnyElement{},
		shouldMatch:      []string{"metric", "group"},
		shouldNotMatch:   []string{"*", "**", "", "{memo2}", "[group={id1}]", "[group1=id]"},
		isFilter:         false,
	},
	{
		namespaceElement: "group1",
		comparableType:   &staticSpecificElement{},
		shouldMatch:      []string{"group1"},
		shouldNotMatch:   []string{"group2", "group", "", "[dyn1=group1]", "[group1]", "*", "**", "{group}"},
		isFilter:         false,
	},
	{
		namespaceElement: "**",
		comparableType:   &staticRecursiveAnyElement{},
		shouldMatch:      []string{"group", "m1", "m2"},
		shouldNotMatch:   []string{"", "*", "**", "{group}"},
		isFilter:         false,
	},
	{
		namespaceElement: "metric1",
		comparableType:   &staticSpecificAcceptingGroupElement{},
		shouldMatch:      []string{"metric1", "[dyn1=metric1]", "[dyn4=metric1]"},
		shouldNotMatch:   []string{"[metric1]", "=metric1]", "[=metric1]", "[dyn1={metric1}]", "*", "**", ""},
		isFilter:         true,
	},
	{
		namespaceElement: "{mem.*[1-3]{1,}}",
		comparableType:   &staticRegexpAcceptingGroupElement{},
		shouldMatch:      []string{"memory3", "mem1", "memo2", "[grp=memo2]"},
		shouldNotMatch:   []string{"memo4", "memory0", "[grp=memory0]", "*", "**", "", "[grp={memo2}]"},
		isFilter:         true,
	},
}

func TestParseNamespaceElement_ValidScenarios(t *testing.T) {
	Convey("Validate ParseNamespace - valid scenarios", t, func() {
		for i, tc := range parseNamespaceElementValidScenarios {
			Convey(fmt.Sprintf("Scenario %d", i), func() {
				// Act
				parsedEl, err := parseNamespaceElement(tc.namespaceElement, tc.isFilter)

				// Assert
				So(err, ShouldBeNil)
				//So(parsedEl.String(), ShouldEqual, tc.namespaceElement)
				So(parsedEl, ShouldHaveSameTypeAs, tc.comparableType)

				// Assert matching (positive)
				for i, m := range tc.shouldMatch {
					Convey(fmt.Sprintf("Scenario %d - Positive matching (%s to %s)", i, m, parsedEl.String()), func() {
						So(parsedEl.Match(m), ShouldBeTrue)
					})
				}

				// Assert matching (negative)
				for i, m := range tc.shouldNotMatch {
					Convey(fmt.Sprintf("Scenario %d - Negative matching (%s to %s)", i, m, parsedEl.String()), func() {
						So(parsedEl.Match(m), ShouldBeFalse)
					})
				}
			})
		}
	})
}

/*****************************************************************************/

func TestParseNamespaceElement_InvalidScenarios(t *testing.T) {
	testCases := []string{
		"",
		"  metr  ",
		" metr",
		"metr ",
		"asd+",
		"{asd[}",
		"[group=]",
		"[=id]",
		"[=]",
		"[gr+]",
		"[group={]",
		"[group={reg[}]",
		"[gr+={id.*}]",
		"[group=id+]",
	}

	Convey("Validate parseNamespaceElement - negative scenarios", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("Scenario %d (%s)", i, tc), func() {
				// Act
				_, err := parseNamespaceElement(tc, false)

				// Assert
				So(err, ShouldBeError)
			})
		}
	})
}
