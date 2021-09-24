//go:build small
// +build small

/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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
