package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMetricDefinitionValidator(t *testing.T) {

	Convey("", t, func() {
		v := NewMetricDefinition()

		// Assert
		So(v.AddRule("/plugin/group1/metric1"), ShouldBeNil)
		So(v.AddRule("/plugin/group2/metric2"), ShouldBeNil)
		So(v.AddRule("/plugin/group2/metric3"), ShouldBeNil)
		//So(v.AddRule("/plugin/group3/[dyn1]/metric4"), ShouldBeNil)
		//So(v.AddRule("/plugin/group4/[dyn2]/metric5"), ShouldBeNil)
		//
		//So(v.AddRule("/plugin/group5/[dyn3]/metric4"), ShouldBeNil) // ok - last element may be repeated if there is no ambiguity
		//So(v.AddRule("/plugin/group6/metric1"), ShouldBeNil)        // ok - last element may be repeated if there is no ambiguity
		//
		//So(v.AddRule("/plugin/[dyn3]/metric6"), ShouldBeError)        // dynamic element on the level where static element is already defined
		//So(v.AddRule("/plugin/group3/[dyn4]/metric7"), ShouldBeError) // 2 dynamic elements on the same level
	})

}
