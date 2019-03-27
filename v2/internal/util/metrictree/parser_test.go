package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseNamespaceElement(t *testing.T) {
	Convey("", t, func() {
		{
			// empty regexp - valid
			el := "{}"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticRegexpElement{})
			So(parsedEl.String(), ShouldEqual, "{}")
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
			el := "*"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldHaveSameTypeAs, &staticAnyElement{})

			So(parsedEl.Match("metric"), ShouldBeTrue)
			So(parsedEl.Match("group"), ShouldBeTrue)
			So(parsedEl.Match(""), ShouldBeTrue)
		}
		{
			// wrong regexp
			el := "{asdsad[}"
			parsedEl := ParseNamespaceElement(el)
			So(parsedEl, ShouldBeNil)
		}

	})
}
