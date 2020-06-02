package runner

import (
	"github.com/fullstorydev/grpchan"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollectorInProcess(collector plugin.InProcessCollector, opt *plugin.Options, inProcChan InProcChanWriteOnly) {
	startCollector(types.NewCollector(collector.Name(), collector.Version(), collector), opt, inProcChan)
}

func StartStreamingCollectorInProcess(collector plugin.InProcessStreamingCollector, opt *plugin.Options, inProcChan InProcChanWriteOnly) {
	startCollector(types.NewStreamingCollector(collector.Name(), collector.Version(), collector), opt, inProcChan)
}

func StartPublisherInProcess(publisher plugin.InProcessPublisher, opt *plugin.Options, inProcChan InProcChanWriteOnly) {
	startPublisher(publisher, publisher.Name(), publisher.Version(), opt, inProcChan)
}

type inProcChan struct {
	GRPCChan chan grpchan.Channel
	MetaCh   chan []byte
}

type InProcChanWriteOnly struct {
	GRPCChan chan<- grpchan.Channel
	MetaCh   chan<- []byte
}

func InProcChan() inProcChan {
	return inProcChan{
		GRPCChan: make(chan grpchan.Channel),
		MetaCh:   make(chan []byte),
	}
}

func (ic inProcChan) WriteOnly() InProcChanWriteOnly {
	return InProcChanWriteOnly{
		GRPCChan: ic.GRPCChan,
		MetaCh:   ic.MetaCh,
	}
}
