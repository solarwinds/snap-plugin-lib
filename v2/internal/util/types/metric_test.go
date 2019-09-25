// +build small

package types

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStaticMetric(t *testing.T) {
	Convey("Validate that methods are working correctly", t, func() {
		// Arrange
		mt := Metric{
			Namespace_: []NamespaceElement{
				{
					Value_: "system",
				},
				{
					Value_: "cpu",
				},
				{
					Value_: "usage",
				},
			},
			Value_:       10,
			Tags_:        map[string]string{"cores": "8", "type": "AMD"},
			Unit_:        "%",
			Timestamp_:   time.Now(),
			Description_: "Usage of processor",
		}

		Convey("Namespace API", func() {
			So(mt.NamespaceText(), ShouldEqual, "/system/cpu/usage")

			So(mt.HasNsElement("system"), ShouldBeTrue)
			So(mt.HasNsElement("cpu"), ShouldBeTrue)
			So(mt.HasNsElement("usage"), ShouldBeTrue)
			So(mt.HasNsElement("ystem"), ShouldBeFalse)
			So(mt.HasNsElement("syste"), ShouldBeFalse)
			So(mt.HasNsElement("/"), ShouldBeFalse)

			So(mt.HasNsElementOn("system", 0), ShouldBeTrue)
			So(mt.HasNsElementOn("cpu", 1), ShouldBeTrue)
			So(mt.HasNsElementOn("usage", 2), ShouldBeTrue)

			So(mt.HasNsElementOn("system", 1), ShouldBeFalse)
			So(mt.HasNsElementOn("cpu", 4), ShouldBeFalse)
			So(mt.HasNsElementOn("usage", -1), ShouldBeFalse)

			So(len(mt.Namespace()), ShouldEqual, 3)

			So(mt.Namespace()[0].IsDynamic(), ShouldBeFalse)
			So(mt.Namespace()[1].IsDynamic(), ShouldBeFalse)
			So(mt.Namespace()[2].IsDynamic(), ShouldBeFalse)
		})

		Convey("Metric API", func() {
			So(mt.Value(), ShouldEqual, 10)

			So(mt.HasTagWithKey("cores"), ShouldBeTrue)
			So(mt.HasTagWithKey("type"), ShouldBeTrue)
			So(mt.HasTagWithKey("AMD"), ShouldBeFalse)
			So(mt.HasTagWithKey(""), ShouldBeFalse)

			So(mt.HasTagWithValue("8"), ShouldBeTrue)
			So(mt.HasTagWithValue("AMD"), ShouldBeTrue)
			So(mt.HasTagWithValue("cores"), ShouldBeFalse)
			So(mt.HasTagWithValue(""), ShouldBeFalse)

			So(mt.HasTag("type", "AMD"), ShouldBeTrue)
			So(mt.HasTag("AMD", "type"), ShouldBeFalse)
		})
	})
}
