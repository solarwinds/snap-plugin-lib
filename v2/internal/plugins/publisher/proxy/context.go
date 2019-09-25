package proxy

import (
	"errors"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type pluginContext struct {
	*proxy.Context

	sessionMts []*types.Metric
}

func NewPluginContext(ctxManager *ContextManager, rawConfig []byte) (*pluginContext, error) {
	if ctxManager == nil {
		return nil, errors.New("can't create context without valid context manager")
	}

	baseContext, err := proxy.NewContext(rawConfig)
	if err != nil {
		return nil, err
	}

	return &pluginContext{
		Context: baseContext,
	}, nil
}

func (pc *pluginContext) ListAllMetrics() []plugin.Metric {
	mts := make([]plugin.Metric, 0, len(pc.sessionMts))

	for _, mt := range pc.sessionMts {
		mts = append(mts, mt)
	}

	return mts
}

func (pc *pluginContext) Count() int {
	return len(pc.sessionMts)
}
