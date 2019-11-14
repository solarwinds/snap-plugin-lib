package pluginrpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
)

func serveInfo(ctx context.Context, statsCh chan *stats.Statistics, pprofAddr string) (*InfoResponse, error) {
	var err error
	response := &InfoResponse{}

	select {
	case statistics := <-statsCh:
		if statistics == nil {
			return response, fmt.Errorf("can't gather statistics (statistics server is not running): %v", err)
		}

		response.Info, err = toGRPCInfo(statistics, pprofAddr)
		if err != nil {
			return response, fmt.Errorf("could't convert statistics to GRPC format: %v", err)
		}
	case <-ctx.Done():
		return response, errors.New("won't retrieve statistics - request canceled")
	}

	return response, nil
}
