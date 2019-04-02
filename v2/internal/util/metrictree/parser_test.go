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
	{ // 5
		namespace:           "/plugin/metric",
		usableForDefinition: true,
		usableForAddition:   true,
	},
}

func TestParseNamespace_ValidScenarios(t *testing.T) {
	Convey("Validate ParseNamespace - valid scenarios", t, func() {
		for i, tc := range parseNamespaceValidScenarios {
			Convey(fmt.Sprintf("Scenario %d", i), func() {
				// Act
				ns, err := ParseNamespace(tc.namespace)

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
		"",
	}

	Convey("Validate ParseNamespace - negative scenarios", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("Scenario %d (%s)", i, tc), func() {
				// Act
				_, err := ParseNamespace(tc)

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
}

var parseNamespaceElementValidScenarios = []parseNamespaceElementValidScenario{
	{ // 0
		namespaceElement: "[group]",
		comparableType:   &dynamicAnyElement{},
		shouldMatch:      []string{"[group=id1]", "[group=id3]", "id3"},
		shouldNotMatch:   []string{"[grp=id1]"},
	},
	{ // 1
		namespaceElement: "[group=id1]",
		comparableType:   &dynamicSpecificElement{},
		shouldMatch:      []string{"id1", "[group=id1]"},
		shouldNotMatch:   []string{"id2", "[group=id2]", "[grp=id1]"},
	},
	{ // 2
		namespaceElement: "[group={id.*}]",
		comparableType:   &dynamicRegexpElement{},
		shouldMatch:      []string{"[group=id1]", "[group=id3]", "id1", "id3"},
		shouldNotMatch:   []string{"[group=i1]", "[grp=id1]", "i1"},
	},
	{ // 3
		namespaceElement: "{}", // valid
		comparableType:   &staticRegexpElement{},
		shouldMatch:      []string{},
		shouldNotMatch:   []string{},
	},
	{ // 4
		namespaceElement: "{mem.*[1-3]{1,}}",
		comparableType:   &staticRegexpElement{},
		shouldMatch:      []string{"memory3", "mem1", "memo2"},
		shouldNotMatch:   []string{"memo4", "memory0"},
	},
	{ // 5
		namespaceElement: "*", // valid
		comparableType:   &staticAnyElement{},
		shouldMatch:      []string{"metric", "group", ""},
		shouldNotMatch:   []string{},
	},
	{
		namespaceElement: "group1",
		comparableType:   &staticSpecificElement{},
		shouldMatch:      []string{"group1"},
		shouldNotMatch:   []string{"group2", "group", ""},
	},
}

func TestParseNamespaceElement_ValidScenarios(t *testing.T) {
	Convey("Validate ParseNamespace - valid scenarios", t, func() {
		for i, tc := range parseNamespaceElementValidScenarios {
			Convey(fmt.Sprintf("Scenario %d", i), func() {
				// Act
				parsedEl, err := parseNamespaceElement(tc.namespaceElement)

				// Assert
				So(err, ShouldBeNil)
				So(parsedEl.String(), ShouldEqual, tc.namespaceElement)
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
				_, err := parseNamespaceElement(tc)

				// Assert
				So(err, ShouldBeError)
			})
		}
	})
}
