package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseNamespaceElement(t *testing.T) {
	Convey("", t, func() {
		{
			// dynamic element - any
			el := "[group]"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &dynamicAnyElement{})
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("[group=id1]"), ShouldBeTrue)
			So(parsedEl.Match("[group=id3]"), ShouldBeTrue)
			So(parsedEl.Match("id3"), ShouldBeTrue)

			So(parsedEl.Match("[grp=id1]"), ShouldBeFalse)
		}
		{
			// dynamic element - concrete
			el := "[group=id1]"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &dynamicSpecificElement{})
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("id1"), ShouldBeTrue)
			So(parsedEl.Match("[group=id1]"), ShouldBeTrue)

			So(parsedEl.Match("id2"), ShouldBeFalse)
			So(parsedEl.Match("[group=id2]"), ShouldBeFalse)
			So(parsedEl.Match("[grp=id1]"), ShouldBeFalse)
		}
		{
			// empty regexp - valid
			el := "{}"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticRegexpElement{})
			So(parsedEl.String(), ShouldEqual, el)
		}
		{
			// some regexp
			el := "{mem.*[1-3]{1,}}"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticRegexpElement{})
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
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticAnyElement{})

			So(parsedEl.Match("metric"), ShouldBeTrue)
			So(parsedEl.Match("group"), ShouldBeTrue)
			So(parsedEl.Match(""), ShouldBeTrue)
		}
		{
			// static concrete match
			el := "group1"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticSpecificElement{})
			So(parsedEl.String(), ShouldEqual, el)

			So(parsedEl.Match("group1"), ShouldBeTrue)
			So(parsedEl.Match("group2"), ShouldBeFalse)
			So(parsedEl.Match("group"), ShouldBeFalse)
			So(parsedEl.Match(""), ShouldBeFalse)
		}
		{
			// wrong regexp
			el := "{asdsad[}"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldBeNil)
		}
	})
}
