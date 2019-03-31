package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMetricDefinitionValidator(t *testing.T) {

	Convey("", t, func() {
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
		So(v.IsValid("/plugin/group1/metric1"), ShouldBeTrue)
		So(v.IsValid("/plugin/group2/metric2"), ShouldBeTrue)
		So(v.IsValid("/plugin/group3/[dyn1=id1]/metric4"), ShouldBeTrue)
		So(v.IsValid("/plugin/group3/id2/metric4"), ShouldBeTrue)
		So(v.IsValid("/plugin/group6/metric1"), ShouldBeTrue)

		// Try to validate (filter) incoming metrics - negative scenarios
		//So(v.IsValid("/plugin/group5/[dyn3]/metric4"), ShouldBeFalse) // todo: no value for dynamic element
		So(v.IsValid("/plugin/group1/metric1/"), ShouldBeFalse)
		So(v.IsValid("/plugin/group1/metric2"), ShouldBeFalse)
		So(v.IsValid("/plugin/group1"), ShouldBeFalse)
		So(v.IsValid("/plugin"), ShouldBeFalse)
		So(v.IsValid("/plugin/[group1=group1]/metric1"), ShouldBeFalse)
		So(v.IsValid("/plugin/group1/metric2/metric2"), ShouldBeFalse)
		So(v.IsValid(""), ShouldBeFalse)
		So(v.IsValid("/"), ShouldBeFalse)
		So(v.IsValid("a/"), ShouldBeFalse)

		// Try to add invalid rules (in current validator state)
		So(v.AddRule("/plugin/[dyn3]/metric6"), ShouldBeError)        // dynamic element on the level where static element is already defined
		So(v.AddRule("/plugin/group3/[dyn4]/metric7"), ShouldBeError) // 2 dynamic elements on the same level
		So(v.AddRule("/plugin/group1/metric1"), ShouldBeError)        // the rules already exists
		So(v.AddRule("/plugin/group3/[dyn1]/metric4"), ShouldBeError) // the rules already exists
	})
}
