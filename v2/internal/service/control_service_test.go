//go:build medium
// +build medium

/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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
	"runtime"
	"testing"
	"time"
)

const (
	monitorTestTimeout  = 3 * time.Second
	memoryLeakTestDelay = 1 * time.Second
)

func routineChecker(t *testing.T, initRoutinesNo int) {
	time.Sleep(memoryLeakTestDelay)

	completeRoutinesNo := runtime.NumGoroutine()
	if initRoutinesNo != completeRoutinesNo {
		t.Fatalf("memory leak (%d != %d)", initRoutinesNo, completeRoutinesNo)
	}
}

func TestControlServiceMonitor_MissingPing(t *testing.T) {
	defer routineChecker(t, runtime.NumGoroutine())

	closeCh := make(chan error)
	doneTestCh := make(chan bool)

	ctx, cancelFn := context.WithCancel(context.Background())
	cs := newControlService(ctx, closeCh, 200*time.Millisecond, 3)

	go func() {
		// ok
		time.Sleep(100 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok, missed 2/3 ping
		time.Sleep(500 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok
		time.Sleep(100 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok, missed 2/3 ping
		time.Sleep(500 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok, missed 2/3 ping
		time.Sleep(500 * time.Millisecond)
		cs.pingCh <- struct{}{}

		doneTestCh <- true
	}()

	select {
	case <-doneTestCh:
		// ok
		cancelFn()
	case <-closeCh:
		t.Fatalf("monitor shouldn't exit")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}

	close(closeCh)
}

func TestControlServiceMonitor_MaxMissedPings(t *testing.T) {
	defer routineChecker(t, runtime.NumGoroutine())

	closeCh := make(chan error)
	doneTestCh := make(chan bool)

	ctx, cancelFn := context.WithCancel(context.Background())
	cs := newControlService(ctx, closeCh, 200*time.Millisecond, 3)

	go func() {
		// ok
		time.Sleep(100 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok, missed 2/3 ping
		time.Sleep(500 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok
		time.Sleep(100 * time.Millisecond)
		cs.pingCh <- struct{}{}

		// ok, missed 3/3 ping
		time.Sleep(700 * time.Millisecond)
		cs.pingCh <- struct{}{} // we should block here

		doneTestCh <- true // shouldn't be executed
	}()

	select {
	case <-closeCh:
		// ok
		<-cs.pingCh // unblock test goroutine
		<-doneTestCh
		cancelFn()
	case <-doneTestCh:
		t.Fatalf("last ping shouldn't have been received")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}

	close(closeCh)
}

func TestControlServiceMonitor_ClosingInfiniteMonitor(t *testing.T) {
	defer routineChecker(t, runtime.NumGoroutine())

	closeCh := make(chan error)
	doneTestCh := make(chan bool)

	ctx, cancelFn := context.WithCancel(context.Background())
	cs := newControlService(ctx, closeCh, 0, 0)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cs.pingCh <- struct{}{}

		time.Sleep(2000 * time.Millisecond)
		cs.pingCh <- struct{}{}

		time.Sleep(100 * time.Millisecond)
		cs.pingCh <- struct{}{}

		doneTestCh <- true
	}()

	select {
	case <-doneTestCh:
		// ok
		cancelFn()
	case <-closeCh:
		t.Fatalf("monitor shouldn't exit")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}

	close(closeCh)
}
