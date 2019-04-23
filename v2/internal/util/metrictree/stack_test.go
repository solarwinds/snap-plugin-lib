package metrictree

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNodeStack(t *testing.T) {
	Convey("Validate that basic operations can be performed on stack", t, func() {
		// Arrange
		s := &nodeStack{}

		// Act
		s.Push(&Node{nodeType: 1})
		s.Push(&Node{nodeType: 2})
		s.Push(&Node{nodeType: 4})

		// Assert
		So(s.Len(), ShouldEqual, 3)

		{
			// Act
			v, ok := s.Pop()
			So(v.nodeType, ShouldEqual, 4)
			So(ok, ShouldBeTrue)
			So(s.Len(), ShouldEqual, 2)
		}
		{
			v, ok := s.Pop()
			So(v.nodeType, ShouldEqual, 2)
			So(ok, ShouldBeTrue)
			So(s.Len(), ShouldEqual, 1)
		}
		{
			v, ok := s.Pop()
			So(v.nodeType, ShouldEqual, 1)
			So(ok, ShouldBeTrue)
			So(s.Len(), ShouldEqual, 0)
		}
		{
			_, ok := s.Pop()
			So(ok, ShouldBeFalse)
			So(s.Len(), ShouldEqual, 0)
		}
	})
}
