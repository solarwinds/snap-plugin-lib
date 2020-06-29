package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
	"github.com/sirupsen/logrus"
)

const (
	DefaultPingTimeout           = 6 * time.Second
	DefaultMaxMissingPingCounter = 3
)

var (
	logControlService = log.WithField("service", "Control")

	RequestedKillError = errors.New("kill requested")
)

type controlService struct {
	pingCh chan struct{}   // notification about received ping
	ctx    context.Context // check for a notification from top level code (service crash etc.)
	errCh  chan error
}

func newControlService(ctx context.Context, errCh chan error, pingTimeout time.Duration, maxMissingPingCounter uint) *controlService {
	cs := &controlService{
		pingCh: make(chan struct{}),
		ctx:    ctx,
		errCh:  errCh,
	}

	go cs.monitor(pingTimeout, maxMissingPingCounter)

	return cs
}

func (cs *controlService) Ping(ctx context.Context, _ *pluginrpc.PingRequest) (*pluginrpc.PingResponse, error) {
	logControlService.Debug("GRPC Ping() received")

	select {
	case <-ctx.Done():
	case cs.pingCh <- struct{}{}:
	}

	return &pluginrpc.PingResponse{}, nil
}

func (cs *controlService) Kill(ctx context.Context, _ *pluginrpc.KillRequest) (*pluginrpc.KillResponse, error) {
	logControlService.Debug("GRPC Kill() received")

	select {
	case <-ctx.Done():
	case cs.errCh <- RequestedKillError:
	}

	return &pluginrpc.KillResponse{}, nil
}

func (cs *controlService) monitor(timeout time.Duration, maxPingMissed uint) {
	pingMissed := uint(0)

	// infinite monitoring (until unload)
	if timeout == time.Duration(0) || maxPingMissed == 0 {
		for {
			select {
			case <-cs.ctx.Done():
				return
			case _, ok := <-cs.pingCh:
				if !ok {
					return
				}
			}
		}
	}

	// monitor for max ping missed
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
				cs.errCh <- fmt.Errorf("ping message missed %d times (timeout: %s)", maxPingMissed, timeout)
				return
			}
		case <-cs.ctx.Done():
			return
		}
	}
}
