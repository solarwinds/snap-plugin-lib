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

func TestDynamicMetric(t *testing.T) {
	Convey("Validate that methods are working correctly", t, func() {
		// Arrange
		mt := Metric{
			Namespace_: []NamespaceElement{
				{
					Value_: "system",
				},
				{
					Value_: "network",
				},
				{
					Name_:        "interface",
					Value_:       "enp0s3",
					Description_: "Name of network interface",
				},
				{
					Value_: "in_bytes",
				},
			},
			Value_:       10,
			Tags_:        map[string]string{},
			Unit_:        "B",
			Timestamp_:   time.Now(),
			Description_: "Bytes received on given interface",
		}

		Convey("Namespace API", func() {
			So(mt.NamespaceText(), ShouldEqual, "/system/network/[interface=enp0s3]/in_bytes")

			So(mt.HasNsElement("system"), ShouldBeTrue)
			So(mt.HasNsElement("network"), ShouldBeTrue)
			So(mt.HasNsElement("in_bytes"), ShouldBeTrue)
			So(mt.HasNsElement("[interface=enp0s3]"), ShouldBeTrue)
			So(mt.HasNsElement("interface"), ShouldBeFalse)
			So(mt.HasNsElement("/"), ShouldBeFalse)

			So(mt.HasNsElementOn("system", 0), ShouldBeTrue)
			So(mt.HasNsElementOn("network", 1), ShouldBeTrue)
			So(mt.HasNsElementOn("in_bytes", 3), ShouldBeTrue)

			So(mt.HasNsElementOn("system", 1), ShouldBeFalse)
			So(mt.HasNsElementOn("network", 4), ShouldBeFalse)
			So(mt.HasNsElementOn("in_bytes", 5), ShouldBeFalse)

			So(len(mt.Namespace()), ShouldEqual, 4)

			So(mt.Namespace()[2].IsDynamic(), ShouldBeTrue)
			So(mt.Namespace()[0].IsDynamic(), ShouldBeFalse)
			So(mt.Namespace()[1].IsDynamic(), ShouldBeFalse)
			So(mt.Namespace()[3].IsDynamic(), ShouldBeFalse)

			So(mt.Namespace()[2].Value(), ShouldEqual, "enp0s3")
			So(mt.Namespace()[2].Name(), ShouldEqual, "interface")
			So(mt.Namespace()[2].Description(), ShouldEqual, "Name of network interface")
		})

		Convey("Metric API", func() {
			So(mt.Value(), ShouldEqual, 10)

			So(mt.HasTagWithKey("cores"), ShouldBeFalse)
			So(mt.HasTagWithKey("AMD"), ShouldBeFalse)
			So(mt.HasTagWithKey(""), ShouldBeFalse)

			So(mt.HasTagWithValue("8"), ShouldBeFalse)
			So(mt.HasTagWithValue(""), ShouldBeFalse)

			So(mt.HasTag("type", "AMD"), ShouldBeFalse)
			So(mt.HasTag("AMD", "type"), ShouldBeFalse)
			So(mt.HasTag("", ""), ShouldBeFalse)
		})
	})
}
