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
		So(v.AddRule("/plugin/group2/metric2"), ShouldBeNil)
		So(v.AddRule("/plugin/group2/metric3"), ShouldBeNil)
		So(v.AddRule("/plugin/group3/[dyn1]/metric4"), ShouldBeNil)
		So(v.AddRule("/plugin/group4/[dyn2]/metric5"), ShouldBeNil)
		So(v.AddRule("/plugin/group5/[dyn3]/metric4"), ShouldBeNil) // ok - last element may be repeated if there is no ambiguity
		So(v.AddRule("/plugin/group6/metric1"), ShouldBeNil)        // ok - last element may be repeated if there is no ambiguity

		// Double-check that rules were applied
		So(len(v.ListRules()), ShouldEqual, 7)

		// Try to validate (filter) incoming metrics - positive scenarios
		validMetricsToAdd := []string{
			"/plugin/group1/metric1",
			"/plugin/group2/metric2",
			"/plugin/group3/[dyn1=id1]/metric4",
			"/plugin/group3/id2/metric4",
			"/plugin/group6/metric1",
		}

		for _, mt := range validMetricsToAdd {
			ok, _ := v.IsValid(mt)
			So(ok, ShouldBeTrue)
		}

		So(v.IsPartiallyValid("/plugin/group1"), ShouldBeTrue)
		So(v.IsPartiallyValid("/plugin/group2"), ShouldBeTrue)
		So(v.IsPartiallyValid("/plugin/group3/[dyn1=id1]"), ShouldBeTrue) // todo: this is not a valid definition
		So(v.IsPartiallyValid("/plugin/group3/id1"), ShouldBeTrue)
		So(v.IsPartiallyValid("/plugin/group6"), ShouldBeTrue)

		// Try to validate (filter) incoming metrics - negative scenarios
		invalidMetricsToAdd := []string{
			"/plugin/group5/[dyn3]/metric4",
			"/plugin/group1/metric1/",
			"/plugin/group1/metric2",
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

		// Try to add invalid rules (in current validator state)
		So(v.AddRule("/plugin/[dyn3]/metric6"), ShouldBeError)        // dynamic element on the level where static element is already defined
		So(v.AddRule("/plugin/group3/[dyn4]/metric7"), ShouldBeError) // 2 dynamic elements on the same level
		So(v.AddRule("/plugin/group1/metric1"), ShouldBeError)        // the rules already exists
		So(v.AddRule("/plugin/group3/[dyn1]/metric4"), ShouldBeError) // the rules already exists
	})
}

func TestMetricFilterValidator(t *testing.T) {

	Convey("Validate that operations can be done on filtering tree", t, func() {
		v := NewMetricFilter()

		// Add valid rules
		So(v.AddRule("/plugin/group1/metric1"), ShouldBeNil)
		So(v.AddRule("/plugin/{id[234]{1,}}/{.*}"), ShouldBeNil)
		So(v.AddRule("/plugin/[group3={id[234]{1,}}]"), ShouldBeNil)
		So(v.AddRule("/plugin/{.*}/group3/{.*}"), ShouldBeNil)
		So(v.AddRule("/plugin/group4/**"), ShouldBeNil)

		// Double-check that rules were applied
		So(len(v.ListRules()), ShouldEqual, 5)

		// Try to validate (filter) incoming metrics - positive scenarios
		validMetricsToAdd := []string{
			"/plugin/group1/metric1",
			"/plugin/id2/metric4",
			"/plugin/id15/group3/metric3",
			"/plugin/group4/m1",
			"/plugin/group4/m1/m2",
		}

		for _, mt := range validMetricsToAdd {
			ok, _ := v.IsValid(mt)
			So(ok, ShouldBeTrue)
		}

		// Try to validate (filter) incoming metrics - negative scenarios
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

}
