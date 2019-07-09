// +build medium

package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStatistics(t *testing.T) {
	startTime := time.Unix(100000, 0)

	Convey("Validate that statistics calculation is correct", t, func() {
		// Arrange
		sc := NewStatsController("example", "1.2.3")
		sc.Run()

		// Act
		sc.UpdateLoadStat(1, "cfg_1", []string{"filt_1_1", "filt_1_2", "filt_1_3"})
		sc.UpdateCollectStat(1, make([]*types.Metric, 4), true, startTime.Add(1*time.Second), startTime.Add(3*time.Second))
		sc.UpdateCollectStat(1, make([]*types.Metric, 6), true, startTime.Add(4*time.Second), startTime.Add(7*time.Second))
		sc.UpdateCollectStat(1, make([]*types.Metric, 11), true, startTime.Add(8*time.Second), startTime.Add(12*time.Second))
		sc.UpdateUnloadStat(1)

		// Assert
		time.Sleep(100 * time.Millisecond)

		So(sc.stats.PluginInfo.Name, ShouldEqual, "example")
		So(sc.stats.PluginInfo.Version, ShouldEqual, "1.2.3")
		So(sc.stats.PluginInfo.StartTime, ShouldNotBeNil)

		//So(sc.stats.TasksDetails, ShouldContainKey, 1)
		//So(sc.stats.TasksDetails[1].TotalMetrics, ShouldContain, 4)

		So(sc.stats.Tasks.CurrentlyActiveTasks, ShouldEqual, 0)
		So(sc.stats.Tasks.TotalActiveTasks, ShouldEqual, 1)

		// Finalize
		stat := <-sc.RequestStat()
		fmt.Printf("stat=%#v\n", stat)

		sc.Close()
	})
}
