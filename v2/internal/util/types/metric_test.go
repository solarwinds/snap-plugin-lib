/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

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

			So(mt.Namespace().Len(), ShouldEqual, 3)

			So(mt.Namespace().At(0).IsDynamic(), ShouldBeFalse)
			So(mt.Namespace().At(1).IsDynamic(), ShouldBeFalse)
			So(mt.Namespace().At(2).IsDynamic(), ShouldBeFalse)
		})

		Convey("Metric API", func() {
			So(mt.Value(), ShouldEqual, 10)

			So(mt.Tags(), ShouldContainKey, "cores")
			So(mt.Tags(), ShouldContainKey, "type")
			So(mt.Tags(), ShouldNotContainKey, "AMD")
			So(mt.Tags(), ShouldNotContainKey, "")

			So(mt.Tags()["type"], ShouldEqual, "AMD")
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

			So(mt.Namespace().Len(), ShouldEqual, 4)

			So(mt.Namespace().At(2).IsDynamic(), ShouldBeTrue)
			So(mt.Namespace().At(0).IsDynamic(), ShouldBeFalse)
			So(mt.Namespace().At(1).IsDynamic(), ShouldBeFalse)
			So(mt.Namespace().At(3).IsDynamic(), ShouldBeFalse)

			So(mt.Namespace().At(2).Value(), ShouldEqual, "enp0s3")
			So(mt.Namespace().At(2).Name(), ShouldEqual, "interface")
			So(mt.Namespace().At(2).Description(), ShouldEqual, "Name of network interface")
		})

		Convey("Tags API", func() {
			So(mt.Value(), ShouldEqual, 10)

			So(mt.Tags(), ShouldNotContainKey, "cores")
			So(mt.Tags(), ShouldNotContainKey, "AMD")
			So(mt.Tags(), ShouldNotContainKey, "")

			So(mt.Tags()["type"], ShouldNotEqual, "AMD")
			So(mt.Tags()["AMD"], ShouldNotEqual, "type")
		})
	})
}
