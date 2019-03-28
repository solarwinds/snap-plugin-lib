// +build small

package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseNamespace(t *testing.T) {
	Convey("", t, func() {
		{
			nsStr := "/plugin/group1/metric"
			ns, err := ParseNamespace(nsStr)

			So(ns, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(ns.isUsableForDefinition(), ShouldBeTrue)
			So(ns.isUsableForAddition(), ShouldBeTrue)

		}
		{
			nsStr := "/plugin/[group2]/metric"
			ns, err := ParseNamespace(nsStr)

			So(ns, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(ns.isUsableForDefinition(), ShouldBeTrue)
			So(ns.isUsableForAddition(), ShouldBeFalse)
		}
		{
			nsStr := "/plugin/[group2=id]/metric"
			ns, err := ParseNamespace(nsStr)

			So(ns, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(ns.isUsableForDefinition(), ShouldBeFalse)
			So(ns.isUsableForAddition(), ShouldBeTrue)
		}
		{
			nsStr := "/plugin/[group2={id.*}]/metric"
			ns, err := ParseNamespace(nsStr)

			So(ns, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(ns.isUsableForDefinition(), ShouldBeFalse)
			So(ns.isUsableForAddition(), ShouldBeFalse)
		}
		{
			nsStr := "/plugin/{id.*}/metric"
			ns, err := ParseNamespace(nsStr)

			So(ns, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(ns.isUsableForDefinition(), ShouldBeFalse)
			So(ns.isUsableForAddition(), ShouldBeFalse)
		}
		{
			nsStr := "/plugin/*/metric"
			ns, err := ParseNamespace(nsStr)

			So(ns, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(ns.isUsableForDefinition(), ShouldBeFalse)
			So(ns.isUsableForAddition(), ShouldBeFalse)
		}
		{
			nsStr := "/"
			_, err := ParseNamespace(nsStr)
			So(err, ShouldBeError)
		}
		{
			nsStr := "sdd/"
			_, err := ParseNamespace(nsStr)
			So(err, ShouldBeError)
		}
		{
			nsStr := "/sdd/sd//df"
			_, err := ParseNamespace(nsStr)
			So(err, ShouldBeError)
		}
		{
			nsStr := ""
			_, err := ParseNamespace(nsStr)
			So(err, ShouldBeError)
		}
	})
}

func TestParseNamespaceElement(t *testing.T) {
	Convey("", t, func() {
		{
			// dynamic element - any
			el := "[group]"
			parsedEl, err := ParseNamespaceElement(el)
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
			parsedEl, err := ParseNamespaceElement(el)
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
			parsedEl, err := ParseNamespaceElement(el)
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
			parsedEl, err := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticRegexpElement{})
			So(err, ShouldBeNil)
			So(parsedEl.String(), ShouldEqual, el)
		}
		{
			// some regexp
			el := "{mem.*[1-3]{1,}}"
			parsedEl, err := ParseNamespaceElement(el)
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
			parsedEl, err := ParseNamespaceElement(el)
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
			parsedEl, err := ParseNamespaceElement(el)
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
			_, err := ParseNamespaceElement(el)
			So(err, ShouldBeError)
		}
	})
}
