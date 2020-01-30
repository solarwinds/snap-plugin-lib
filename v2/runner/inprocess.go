package runner

import (
	"github.com/fullstorydev/grpchan"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollectorInProcess(collector plugin.InProcessCollector, name string, version string, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	startCollector(collector, name, version, opt, grpcChan)
}

func StartPublisherInProcess(publisher plugin.InProcessPublisher, name string, version string, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	startPublisher(publisher, name, version, opt, grpcChan)
}
