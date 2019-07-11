// +build medium

package stats

import (
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	pluginName    = "example"
	pluginVersion = "1.2.3"
)

func TestStatistics(t *testing.T) {
	startTime := time.Unix(100000, 0)

	Convey("Validate that statistics calculation is correct", t, func() {
		// Arrange
		sc := NewStatsController(pluginName, pluginVersion, &types.Options{}).(*StatisticsController)

		// Act
		sc.UpdateLoadStat(1, "cfg_1", []string{"filt_1_1", "filt_1_2", "filt_1_3"})
		sc.UpdateCollectStat(1, make([]*types.Metric, 4), true, startTime.Add(1*time.Second), startTime.Add(3*time.Second))
		sc.UpdateCollectStat(1, make([]*types.Metric, 6), true, startTime.Add(4*time.Second), startTime.Add(7*time.Second))
		sc.UpdateCollectStat(1, make([]*types.Metric, 11), true, startTime.Add(8*time.Second), startTime.Add(12*time.Second))

		// Assert
		time.Sleep(100 * time.Millisecond)

		pi := &sc.stats.PluginInfo

		So(pi.Name, ShouldEqual, pluginName)
		So(pi.Version, ShouldEqual, pluginVersion)
		So(pi.Started.Time, ShouldNotBeNil)

		So(sc.stats.TasksDetails, ShouldContainKey, 1)
		So(sc.stats.TasksDetails[1].Counters.TotalMetrics, ShouldEqual, 21)

		So(sc.stats.TasksSummary.Counters.CurrentlyActiveTasks, ShouldEqual, 1)
		So(sc.stats.TasksSummary.Counters.TotalActiveTasks, ShouldEqual, 1)

		// Finalize
		sc.Close()
	})
}
