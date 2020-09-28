/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
	"github.com/solarwinds/snap-plugin-lib/v2/pluginrpc"
)

const (
	DefaultPingTimeout           = 6 * time.Second
	DefaultMaxMissingPingCounter = 3
)

var (
	controlSrvFields = logrus.Fields{"service": "Control"}

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
	cs.logger().WithFields(controlSrvFields).Debug("GRPC Ping() received")

	select {
	case <-ctx.Done():
	case cs.pingCh <- struct{}{}:
	}

	return &pluginrpc.PingResponse{}, nil
}

func (cs *controlService) Kill(ctx context.Context, _ *pluginrpc.KillRequest) (*pluginrpc.KillResponse, error) {
	cs.logger().WithFields(controlSrvFields).Debug("GRPC Kill() received")

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
			cs.logger().WithFields(controlSrvFields).WithFields(logrus.Fields{
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

func (cs *controlService) logger() logrus.FieldLogger {
	return log.WithCtx(cs.ctx).WithFields(moduleFields).WithField("service", "Control")
}
