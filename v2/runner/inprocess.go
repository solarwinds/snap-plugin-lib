package runner

import (
	"github.com/fullstorydev/grpchan"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollectorInProcess(collector plugin.InProcessCollector, opt *plugin.Options, grpcChan chan<- grpchan.Channel, metaCh chan<- []byte) {
	startCollector(types.NewCollector(collector.Name(), collector.Version(), collector), opt, grpcChan, metaCh)
}

func StartStreamingCollectorInProcess(collector plugin.InProcessStreamingCollector, opt *plugin.Options, grpChan chan<- grpchan.Channel, metaCh chan<- []byte) {
	startCollector(types.NewStreamingCollector(collector.Name(), collector.Version(), collector), opt, grpChan, metaCh)
}

func StartPublisherInProcess(publisher plugin.InProcessPublisher, opt *plugin.Options, grpcChan chan<- grpchan.Channel, metaCh chan<- []byte) {
	startPublisher(publisher, publisher.Name(), publisher.Version(), opt, grpcChan, metaCh)
}
