package pluginrpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	DefaultPingTimeout           = 3 * time.Second
	DefaultMaxMissingPingCounter = 3
)

var (
	logControlService = log.WithField("service", "Control")

	RequestedKillError = errors.New("kill requested")
)

type controlService struct {
	pingCh  chan struct{} // notification about received ping
	closeCh chan error    // request exit to main routine
}

func newControlService(closeCh chan error, pingTimeout time.Duration, maxMissingPingCounter uint) *controlService {
	cs := &controlService{
		pingCh:  make(chan struct{}),
		closeCh: closeCh,
	}

	if pingTimeout != time.Duration(0) && maxMissingPingCounter != 0 {
		go cs.monitor(pingTimeout, maxMissingPingCounter)
	} else {
		go func() {
			for {
				_, ok := <-cs.pingCh
				if !ok {
					return
				}
			}
		}()
	}

	return cs
}

func (cs *controlService) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	logControlService.Debug("GRPC Ping() received")

	cs.pingCh <- struct{}{}

	return &PingResponse{}, nil
}

func (cs *controlService) Kill(context.Context, *KillRequest) (*KillResponse, error) {
	logControlService.Debug("GRPC Kill() received")

	cs.closeCh <- RequestedKillError

	return &KillResponse{}, nil
}

func (cs *controlService) monitor(timeout time.Duration, maxPingMissed uint) {
	pingMissed := uint(0)

	for {
		select {
		case <-cs.pingCh:
			pingMissed = 0
		case <-time.After(timeout):
			pingMissed++
			log.WithFields(logrus.Fields{
				"missed": pingMissed,
				"max":    maxPingMissed,
			}).Warningf("Ping timeout occurred")

			if pingMissed >= maxPingMissed {
				cs.closeCh <- fmt.Errorf("ping message missed %d times (timeout: %s)", maxPingMissed, timeout)
				return
			}
		}
	}
}
