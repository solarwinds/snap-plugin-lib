// +build small

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

package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMetricDefinitionValidator(t *testing.T) {
	Convey("Validate that operations can be done on definition tree", t, func() {
		v := NewMetricDefinition()

		// Add valid rules
		So(v.AddRule("/plugin/group1/metric1"), ShouldBeNil)
		So(v.AddRule(".plugin.group2.metric2"), ShouldBeNil)
		So(v.AddRule("/plugin/group2/metric3"), ShouldBeNil)
		So(v.AddRule("/plugin/group3/[dyn1]/metric4"), ShouldBeNil)
		So(v.AddRule("/plugin/group4/[dyn2]/metric5"), ShouldBeNil)
		So(v.AddRule("@plugin@group5@[dyn3]@metric4"), ShouldBeNil) // ok - last element may be repeated if there is no ambiguity
		So(v.AddRule("/plugin/group6/metric1"), ShouldBeNil)        // ok - last element may be repeated if there is no ambiguity

		// Double-check that rules were applied
		defer func() {
			So(len(v.ListRules()), ShouldEqual, 7)
		}()

		Convey("Try to validate (filter) incoming metrics - positive scenarios", func() {
			validMetricsToAdd := []string{
				"/plugin/group1/metric1",
				"/plugin/group2/metric2",
				".plugin.group3.[dyn1=/id1/].metric4",
				"/plugin/group3/id2/metric4",
				"/plugin/group6/metric1",
				"/plugin/group5/[dyn3]/metric4",
			}

			for _, mt := range validMetricsToAdd {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeTrue)
			}

			So(v.IsPartiallyValid("/plugin/group1"), ShouldBeTrue)
			So(v.IsPartiallyValid("/plugin/group2"), ShouldBeTrue)
			So(v.IsPartiallyValid("-plugin-group3-[dyn1=/id1]"), ShouldBeTrue)
			So(v.IsPartiallyValid("/plugin/group3/id1"), ShouldBeTrue)
			So(v.IsPartiallyValid("/plugin/group6"), ShouldBeTrue)
		})

		Convey("Try to validate (filter) incoming metrics - negative scenarios", func() {
			invalidMetricsToAdd := []string{
				"/plugin/group1/metric1/",
				".plugin.group1.metric2",
				"/plugin/group1",
				"/plugin",
				"/plugin/[group1=group1]/metric1",
				"/plugin/group1/metric2/metric2",
				"",
				"/",
				"a/",
			}

			for _, mt := range invalidMetricsToAdd {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeFalse)
			}
		})

		Convey("Try to add invalid rules (in current validator state)", func() {
			So(v.AddRule("/plugin/[dyn3]/metric6"), ShouldBeError)        // dynamic element on the level where static element is already defined
			So(v.AddRule("/plugin/group3/[dyn4]/metric7"), ShouldBeError) // 2 dynamic elements on the same level
			So(v.AddRule("/plugin/group1/metric1"), ShouldBeError)        // the rules already exists
			So(v.AddRule("/plugin/group3/[dyn1]/metric4"), ShouldBeError) // the rules already exists
		})

		Convey("Try some edge cases to validate default constraints", func() {
			So(v.AddRule("/plugin/group7/metric8/[dyn4]"), ShouldBeError) //  can't define dynamic leaf
			So(v.AddRule("/plugin/group1/metric1/[dyn4]"), ShouldBeError) // can't define dynamic leaf

			for _, mt := range []string{
				"/plugin/group7/metric8/another1", // can't submit a metric that was not defined previously
				"/plugin/group3/[dyn1]/metric6",   // can't submit a metric that was not defined previously
				"/plugin/group1/metric1/another1", // can't redefine element as non-leaf one
			} {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeFalse)
			}
		})
	})
}

func TestMetricFilterValidator_NoDefinition(t *testing.T) {
	Convey("Validate that operations can be done on filtering tree", t, func() {
		d := NewMetricDefinition()
		v := NewMetricFilter(d)

		// Add valid rules
		So(v.AddRule("/plugin/group1/metric1"), ShouldBeNil)
		So(v.AddRule("_plugin_{id[234]{1,}}_{.*}"), ShouldBeNil)
		So(v.AddRule("/plugin/{.*}/group3/{.*}"), ShouldBeNil)
		So(v.AddRule("/plugin/group4/**"), ShouldBeNil)

		// Double-check that rules were applied
		defer func() {
			So(len(v.ListRules()), ShouldEqual, 4)
		}()

		Convey("Add invalid rules", func() {
			So(v.AddRule("/plugin/[group3={id[234]{1,}}]"), ShouldBeError) // dynamic element with no definition
			So(v.AddRule("/plugin"), ShouldBeError)                        // len < 2
			So(v.AddRule("/plugin/{af[}/metric4"), ShouldBeError)          // invalid regexp
		})

		Convey("Try to validate (filter) incoming metrics - positive scenarios", func() {
			validMetricsToAdd := []string{
				"/plugin/group1/metric1",
				"/plugin/id2/metric4",
				".plugin.id15.group3.metric3",
				"/plugin/group4/m1",
				"/plugin/group4/m1/m2",
			}

			for _, mt := range validMetricsToAdd {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeTrue)
			}
		})

		Convey("Try to validate (filter) incoming metrics - negative scenarios", func() {
			invalidMetricsToAdd := []string{
				"/plugin/group2/metric4",
				"/plugin/[group2=group2]/metric4",
				"/plugin/id15/group4/metric4",
				"/plugin/group4",
			}

			for _, mt := range invalidMetricsToAdd {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeFalse)
			}
		})
	})

}

func TestMetricFilterValidator_MetricDefinition(t *testing.T) {
	Convey("Validate that operations can be done on filtering tree", t, func() {
		d := NewMetricDefinition()
		v := NewMetricFilter(d)

		// Add valid definition rules
		So(d.AddRule("/plugin/group1/[dyn1]/metric1"), ShouldBeNil)
		So(d.AddRule("/plugin/group2/sub1/metric1"), ShouldBeNil)
		So(d.AddRule(".plugin.group2.sub2.metric2"), ShouldBeNil)
		So(d.AddRule("/plugin/group2/sub3/metric3"), ShouldBeNil)
		So(d.AddRule("_plugin_group3_[dyn2]_[dyn3]_metric2"), ShouldBeNil)
		So(d.AddRule("/plugin/group4/[dyn1]/[dyn3]/metric2"), ShouldBeNil)

		// Add valid filtering rules
		So(v.AddRule("/plugin/group1/id1/metric1"), ShouldBeNil)
		So(v.AddRule("/plugin/group2/*/{metric[123]+}"), ShouldBeNil)
		So(v.AddRule("/plugin/group3/id1/[dyn3=id2]/metric2"), ShouldBeNil)
		So(v.AddRule("/plugin/group3/{id3+}/[dyn3={id4+}]/metric2"), ShouldBeNil)
		So(v.AddRule("_plugin_group4_**"), ShouldBeNil)

		// Double-check that rules were applied
		defer func() {
			So(len(d.ListRules()), ShouldEqual, 6)
			So(len(v.ListRules()), ShouldEqual, 5)
		}()

		Convey("Add invalid filtering rules (no compatible with definitions)", func() {
			So(v.AddRule(".plugin.group1.id1.metric12"), ShouldBeError)
			So(v.AddRule("/plugins/group1/id1/metric1"), ShouldBeError)
			So(v.AddRule("/plugin/group2/[sub2]/{metric[123]+}"), ShouldBeError)
			So(v.AddRule("/plugin/group3/[dyn2]/[dyn4=val]/metric2"), ShouldBeError)
		})

		Convey("Try to validate (filter) incoming metrics - positive scenarios", func() {
			validMetricsToAdd := []string{
				"/plugin/group1/id1/metric1",
				"/plugin/group2/sub2/metric2",
				"/plugin/group3/id1/[dyn3=id2]/metric2",
				".plugin.group3.[dyn2=id1].[dyn3=id2].metric2",
			}

			for _, mt := range validMetricsToAdd {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeTrue)
			}
		})

		Convey("Try to validate (filter) incoming metrics - negative scenarios", func() {
			invalidMetricsToAdd := []string{
				"/plugin/group1/id1/metric4",
				"/plugin/group2/sub2/metric4",
				"/plugin/group3/[dyn2=id2]/[dyn3=id2]/metric2",
			}

			for _, mt := range invalidMetricsToAdd {
				ok, _ := v.IsValid(mt)
				So(ok, ShouldBeFalse)
			}
		})
	})

}
