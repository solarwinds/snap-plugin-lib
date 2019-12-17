// +build medium

package pluginrpc

import (
	"testing"
	"time"
)

const (
	monitorTestTimeout = 3 * time.Second
)

func TestControlServiceMonitor_MissingPing(t *testing.T) {
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
}

func TestControlServiceMonitor_MaxMissedPings(t *testing.T) {
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
		// ok
	case <-doneTestCh:
		t.Fatalf("last ping shouldn't have been received")
	case <-time.After(monitorTestTimeout):
		t.Fatalf("test timeout")
	}
}
