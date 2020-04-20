// +build medium

package service

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

const (
	monitorTestTimeout  = 3 * time.Second
	memoryLeakTestDelay = 1 * time.Second
)

func TestControlServiceMonitor_MissingPing(t *testing.T) {
	initGoroutines := runtime.NumGoroutine()

	closeCh := make(chan error)
	doneTestCh := make(chan bool)

	cs := newControlService(closeCh, 200*time.Millisecond, 3)

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
	case <-closeCh:
		t.Fatalf("monitor shouldn't exit")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}

	close(closeCh)
	time.Sleep(memoryLeakTestDelay)

	if initGoroutines != runtime.NumGoroutine() {
		t.Fatalf("memory leak")
	}
}

func TestControlServiceMonitor_MaxMissedPings(t *testing.T) {
	initGoroutines := runtime.NumGoroutine()

	closeCh := make(chan error)
	doneTestCh := make(chan bool)

	cs := newControlService(closeCh, 200*time.Millisecond, 3)

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
		// ok, unblock test routine to avoid a leak
		<-cs.pingCh
		<-doneTestCh
	case <-doneTestCh:
		t.Fatalf("last ping shouldn't have been received")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}

	close(closeCh)
	time.Sleep(memoryLeakTestDelay)

	if initGoroutines != runtime.NumGoroutine() {
		fmt.Print(runtime.NumGoroutine())
		t.Fatalf("memory leak")
	}
}

func TestControlServiceMonitor_ClosingInfinitiveMonitor(t *testing.T) {
	initGoroutines := runtime.NumGoroutine()

	closeCh := make(chan error)
	doneTestCh := make(chan bool)

	cs := newControlService(closeCh, 0, 0)

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
	case <-closeCh:
		t.Fatalf("monitor shouldn't exit")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}

	close(closeCh)
	time.Sleep(memoryLeakTestDelay)

	if initGoroutines != runtime.NumGoroutine() {
		t.Fatalf("memory leak")
	}
}
