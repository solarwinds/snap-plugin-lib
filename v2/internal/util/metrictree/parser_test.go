// +build small

package metrictree

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*****************************************************************************/

type parseNamespaceValidScenario struct {
	namespace           string
	usableForDefinition bool
	usableForAddition   bool
}

var parseNamespaceValidScenarios = []parseNamespaceValidScenario{
	{ // 0
		namespace:           "/plugin/group1/metric",
		usableForDefinition: true,
		usableForAddition:   true,
	},
	{ // 1
		namespace:           "/plugin/[group2]/metric",
		usableForDefinition: true,
		usableForAddition:   false,
	},
	{ // 2
		namespace:           "/plugin/[group2=id]/metric",
		usableForDefinition: false,
		usableForAddition:   true,
	},
	{ // 3
		namespace:           "/plugin/[group2={id.*}]/metric",
		usableForDefinition: false,
		usableForAddition:   false,
	},
	{ // 4
		namespace:           "/plugin/{id.*}/metric",
		usableForDefinition: false,
		usableForAddition:   false,
	},
	{ // 5
		namespace:           "/plugin/*/metric",
		usableForDefinition: false,
		usableForAddition:   false,
	},
	{ // 6
		namespace:           "/plugin/metric/**",
		usableForDefinition: false,
		usableForAddition:   false,
	},
	{ // 7
		namespace:           "/plugin/metric",
		usableForDefinition: true,
		usableForAddition:   true,
	},
	{ // 7
		namespace:           "/plugin/**",
		usableForDefinition: false,
		usableForAddition:   false,
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
				So(ns.isUsableForDefinition(), ShouldEqual, tc.usableForDefinition)
				So(ns.isUsableForAddition(), ShouldEqual, tc.usableForAddition)
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
		shouldMatch:      []string{"[group=id1]", "[group=id3]", "id3"},
		shouldNotMatch:   []string{"[grp=id1]", "[group]"},
		isFilter:         false,
	},
	{
		namespaceElement: "[group=id1]",
		comparableType:   &dynamicSpecificElement{},
		shouldMatch:      []string{"id1", "[group=id1]"},
		shouldNotMatch:   []string{"id2", "[group=id2]", "[grp=id1]"},
		isFilter:         false,
	},
	{
		namespaceElement: "[group={id.*}]",
		comparableType:   &dynamicRegexpElement{},
		shouldMatch:      []string{"[group=id1]", "[group=id3]", "id1", "id3"},
		shouldNotMatch:   []string{"[group=i1]", "[grp=id1]", "i1"},
		isFilter:         false,
	},
	{
		namespaceElement: "{}", // valid
		comparableType:   &staticRegexpElement{},
		shouldMatch:      []string{},
		shouldNotMatch:   []string{},
		isFilter:         false,
	},
	{
		namespaceElement: "{mem.*[1-3]{1,}}",
		comparableType:   &staticRegexpElement{},
		shouldMatch:      []string{"memory3", "mem1", "memo2"},
		shouldNotMatch:   []string{"memo4", "memory0", "group"},
		isFilter:         false,
	},
	{
		namespaceElement: "*", // valid
		comparableType:   &staticAnyElement{},
		shouldMatch:      []string{"metric", "group", ""},
		shouldNotMatch:   []string{},
		isFilter:         false,
	},
	{
		namespaceElement: "group1",
		comparableType:   &staticSpecificElement{},
		shouldMatch:      []string{"group1"},
		shouldNotMatch:   []string{"group2", "group", "", "[dyn1=group1]"},
		isFilter:         false,
	},
	{
		namespaceElement: "**",
		comparableType:   &staticRecursiveAnyElement{},
		shouldMatch:      []string{"group", "m1", "m2"},
		shouldNotMatch:   []string{},
		isFilter:         false,
	},
	{
		/* special case:
		when we have metric defined: /plugin/[dyn1]/metric1
		we can add filter using two methods:
			/plugin/id1/metric1
			/plugin/[dyn1=id]/metric1
		we need to be able to add metric using form:
			/plugin/id1/metric1
		*/
		namespaceElement: "metric1",
		comparableType:   &staticSpecificAcceptingGroupElement{},
		shouldMatch:      []string{"metric1", "[dyn1=metric1]", "[dyn4=metric1]"},
		shouldNotMatch:   []string{"[metric1]", "=metric1]", "[=metric1]"},
		isFilter:         true,
	},
	//{
	//	namespaceElement: "{mem.*[1-3]{1,}}",
	//	comparableType:   &staticRegexpElement{},
	//	shouldMatch:      []string{"memory3", "mem1", "memo2"},
	//	shouldNotMatch:   []string{"memo4", "memory0"},
	//	isFilter:         true,
	//},

}

func TestParseNamespaceElement_ValidScenarios(t *testing.T) {
	Convey("Validate ParseNamespace - valid scenarios", t, func() {
		for i, tc := range parseNamespaceElementValidScenarios[8:] {
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
						fmt.Printf("***%s->%s \n", tc.namespaceElement, m)
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
