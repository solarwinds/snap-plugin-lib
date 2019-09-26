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
			So(mt.Namespace().String(), ShouldEqual, "/system/cpu/usage")

			So(mt.Namespace().HasElement("system"), ShouldBeTrue)
			So(mt.Namespace().HasElement("cpu"), ShouldBeTrue)
			So(mt.Namespace().HasElement("usage"), ShouldBeTrue)
			So(mt.Namespace().HasElement("ystem"), ShouldBeFalse)
			So(mt.Namespace().HasElement("syste"), ShouldBeFalse)
			So(mt.Namespace().HasElement("/"), ShouldBeFalse)

			So(mt.Namespace().HasElementOn("system", 0), ShouldBeTrue)
			So(mt.Namespace().HasElementOn("cpu", 1), ShouldBeTrue)
			So(mt.Namespace().HasElementOn("usage", 2), ShouldBeTrue)

			So(mt.Namespace().HasElementOn("system", 1), ShouldBeFalse)
			So(mt.Namespace().HasElementOn("cpu", 4), ShouldBeFalse)
			So(mt.Namespace().HasElementOn("usage", -1), ShouldBeFalse)

			//So(mt.Namespace()[0].IsDynamic(), ShouldBeFalse) // todo: uncomment
			//So(mt.Namespace()[1].IsDynamic(), ShouldBeFalse)
			//So(mt.Namespace()[2].IsDynamic(), ShouldBeFalse)
		})

		Convey("Metric API", func() {
			So(mt.Value(), ShouldEqual, 10)

			So(mt.Tags().ContainsKey("cores"), ShouldBeTrue)
			So(mt.Tags().ContainsKey("type"), ShouldBeTrue)
			So(mt.Tags().ContainsKey("AMD"), ShouldBeFalse)
			So(mt.Tags().ContainsKey(""), ShouldBeFalse)

			So(mt.Tags().ContainsValue("8"), ShouldBeTrue)
			So(mt.Tags().ContainsValue("AMD"), ShouldBeTrue)
			So(mt.Tags().ContainsValue("cores"), ShouldBeFalse)
			So(mt.Tags().ContainsValue(""), ShouldBeFalse)

			So(mt.Tags().Contains("type", "AMD"), ShouldBeTrue)
			So(mt.Tags().Contains("AMD", "type"), ShouldBeFalse)
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
			So(mt.Namespace().String(), ShouldEqual, "/system/network/[interface=enp0s3]/in_bytes")

			So(mt.Namespace().HasElement("system"), ShouldBeTrue)
			So(mt.Namespace().HasElement("network"), ShouldBeTrue)
			So(mt.Namespace().HasElement("in_bytes"), ShouldBeTrue)
			So(mt.Namespace().HasElement("[interface=enp0s3]"), ShouldBeTrue)
			So(mt.Namespace().HasElement("interface"), ShouldBeFalse)
			So(mt.Namespace().HasElement("/"), ShouldBeFalse)

			So(mt.Namespace().HasElementOn("system", 0), ShouldBeTrue)
			So(mt.Namespace().HasElementOn("network", 1), ShouldBeTrue)
			So(mt.Namespace().HasElementOn("in_bytes", 3), ShouldBeTrue)

			So(mt.Namespace().HasElementOn("system", 1), ShouldBeFalse)
			So(mt.Namespace().HasElementOn("network", 4), ShouldBeFalse)
			So(mt.Namespace().HasElementOn("in_bytes", 5), ShouldBeFalse)

			//So(len(mt.Namespace()), ShouldEqual, 4) // todo: uncomment

			//So(mt.Namespace()[2].IsDynamic(), ShouldBeTrue)
			//So(mt.Namespace()[0].IsDynamic(), ShouldBeFalse)
			//So(mt.Namespace()[1].IsDynamic(), ShouldBeFalse)
			//So(mt.Namespace()[3].IsDynamic(), ShouldBeFalse)
			//
			//So(mt.Namespace()[2].Value(), ShouldEqual, "enp0s3")
			//So(mt.Namespace()[2].Name(), ShouldEqual, "interface")
			//So(mt.Namespace()[2].Description(), ShouldEqual, "Name of network interface")
		})

		Convey("Tags API", func() {
			So(mt.Value(), ShouldEqual, 10)

			So(mt.Tags().ContainsKey("cores"), ShouldBeFalse)
			So(mt.Tags().ContainsKey("AMD"), ShouldBeFalse)
			So(mt.Tags().ContainsKey(""), ShouldBeFalse)

			So(mt.Tags().ContainsValue("8"), ShouldBeFalse)
			So(mt.Tags().ContainsValue(""), ShouldBeFalse)

			So(mt.Tags().Contains("type", "AMD"), ShouldBeFalse)
			So(mt.Tags().Contains("AMD", "type"), ShouldBeFalse)
			So(mt.Tags().Contains("", ""), ShouldBeFalse)
		})
	})
}
