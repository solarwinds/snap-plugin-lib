package runner

import (
	"github.com/fullstorydev/grpchan"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollectorInProcess(collector plugin.InProcessCollector, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	startCollector(types.NewCollector(collector.Name(), collector.Version(), collector), opt, grpcChan)
}

func StartPublisherInProcess(publisher plugin.InProcessPublisher, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	startPublisher(publisher, publisher.Name(), publisher.Version(), opt, grpcChan)
}
