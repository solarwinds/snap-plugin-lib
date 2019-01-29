// +build small

package plugin

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type splitTestCase struct {
	length        int
	partialLength int
	expectedSplit []Range
}

var splitTestCases = []splitTestCase{
	{ // 0
		length:        50,
		partialLength: 20,
		expectedSplit: []Range{{0, 20}, {20, 40}, {40, 50}},
	},
	{ // 1
		length:        40,
		partialLength: 20,
		expectedSplit: []Range{{0, 20}, {20, 40}},
	},
	{ // 2
		length:        39,
		partialLength: 20,
		expectedSplit: []Range{{0, 20}, {20, 39}},
	},
	{ // 3
		length:        41,
		partialLength: 20,
		expectedSplit: []Range{{0, 20}, {20, 40}, {40, 41}},
	},
	{ // 4
		length:        20,
		partialLength: 30,
		expectedSplit: []Range{{0, 20}},
	},
	{ // 5
		length:        0,
		partialLength: 10,
		expectedSplit: []Range(nil),
	},
	{ // 6
		length:        0,
		partialLength: 0,
		expectedSplit: []Range(nil),
	},
	{ // 7
		length:        -1,
		partialLength: -2,
		expectedSplit: []Range(nil),
	},
	{ // 8
		length:        1,
		partialLength: -1,
		expectedSplit: []Range(nil),
	},
	{ // 9
		length:        -1,
		partialLength: 2,
		expectedSplit: []Range(nil),
	},
}

func TestCalculateChunkIndexes(t *testing.T) {
	Convey("Validate CalculateChunkIndexes function", t, func() {
		for id, testCase := range splitTestCases {
			Convey(fmt.Sprintf("Scenario %d", id), func() {
				// Act
				result := CalculateChunkIndexes(testCase.length, testCase.partialLength)

				// Assert
				So(len(result), ShouldEqual, len(testCase.expectedSplit))
				So(result, ShouldResemble, testCase.expectedSplit)
			})
		}
	})
}

func TestChunkMetrics(t *testing.T) {
	Convey("Validate ChunkMetrics function", t, func() {
		// Arrange
		var mts []Metric
		for i := 0; i < 35; i++ {
			mts = append(mts, Metric{
				Namespace: NewNamespace("test", fmt.Sprintf("metric_%d", i)),
				Data:      i*10 + 5,
			})
		}

		// Act
		splitMts := ChunkMetrics(mts, 20)

		// Assert
		So(len(splitMts), ShouldEqual, 2)
		So(len(splitMts[0]), ShouldEqual, 20)
		So(len(splitMts[1]), ShouldEqual, 15)

		So(splitMts[0][0].Namespace.String(), ShouldEqual, "/test/metric_0")
		So(splitMts[0][0].Data, ShouldEqual, 5)

		So(splitMts[0][19].Namespace.String(), ShouldEqual, "/test/metric_19")
		So(splitMts[0][19].Data, ShouldEqual, 195)

		So(splitMts[1][0].Namespace.String(), ShouldEqual, "/test/metric_20")
		So(splitMts[1][0].Data, ShouldEqual, 205)

		So(splitMts[1][1].Namespace.String(), ShouldEqual, "/test/metric_21")
		So(splitMts[1][1].Data, ShouldEqual, 215)

		So(splitMts[1][14].Namespace.String(), ShouldEqual, "/test/metric_34")
		So(splitMts[1][14].Data, ShouldEqual, 345)
	})
}
