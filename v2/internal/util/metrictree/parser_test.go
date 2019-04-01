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

func TestParseNamespace(t *testing.T) {
	testCases := []string{
		"/",
		"el1/",
		"/el1/el2//m3",
		"/el1/el_#/m4",
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

func TestParseNamespaceElement(t *testing.T) {
	Convey("", t, func() {
		{
			// dynamic element - any
			el := "[group]"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &dynamicAnyElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("[group=id1]"), ShouldBeTrue)
			So(parsedEl.Match("[group=id3]"), ShouldBeTrue)
			So(parsedEl.Match("id3"), ShouldBeTrue)

			So(parsedEl.Match("[grp=id1]"), ShouldBeFalse)
		}
		{
			// dynamic element - concrete
			el := "[group=id1]"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &dynamicSpecificElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("id1"), ShouldBeTrue)
			So(parsedEl.Match("[group=id1]"), ShouldBeTrue)

			So(parsedEl.Match("id2"), ShouldBeFalse)
			So(parsedEl.Match("[group=id2]"), ShouldBeFalse)
			So(parsedEl.Match("[grp=id1]"), ShouldBeFalse)
		}
		{
			el := "[group={id.*}]"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &dynamicRegexpElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("[group=id1]"), ShouldBeTrue)
			So(parsedEl.Match("[group=id3]"), ShouldBeTrue)
			So(parsedEl.Match("id1"), ShouldBeTrue)
			So(parsedEl.Match("id3"), ShouldBeTrue)

			So(parsedEl.Match("[group=i1]"), ShouldBeFalse)
			So(parsedEl.Match("[grp=id1]"), ShouldBeFalse)
			So(parsedEl.Match("i1"), ShouldBeFalse)
		}
		{
			// empty regexp - valid
			el := "{}"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticRegexpElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)
		}
		{
			// some regexp
			el := "{mem.*[1-3]{1,}}"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticRegexpElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("memory3"), ShouldBeTrue)
			So(parsedEl.Match("mem1"), ShouldBeTrue)
			So(parsedEl.Match("memo2"), ShouldBeTrue)

			So(parsedEl.Match("memo4"), ShouldBeFalse)
			So(parsedEl.Match("memory0"), ShouldBeFalse)
		}
		{
			// static any match
			el := "*"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticAnyElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("metric"), ShouldBeTrue)
			So(parsedEl.Match("group"), ShouldBeTrue)
			So(parsedEl.Match(""), ShouldBeTrue)
		}
		{
			// static concrete match
			el := "group1"
			parsedEl, err := parseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticSpecificElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("group1"), ShouldBeTrue)
			So(parsedEl.Match("group2"), ShouldBeFalse)
			So(parsedEl.Match("group"), ShouldBeFalse)
			So(parsedEl.Match(""), ShouldBeFalse)
		}
		{
			// wrong regexp
			el := "{asdsad[}"
			_, err := parseNamespaceElement(el)
			So(err, ShouldBeError)
		}
	})
}
