package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

/*
#include <stdlib.h>

typedef char * (collectCallbackT)();

typedef struct {
	collectCallbackT *collectCallback;
	//loadCallback loadCallbackT
	//unloadCallback loadCallbackT
	//definePluginCallback definePluginCallbackT
} cCollectorT;

void callCollect(collectCallbackT collectCallback) {
	cCollectorT cc;
	cc.collectCallback = collectCallback;
	cc.collectCallback();
}

*/
import "C"

type bridgeCollector struct {
	cCollector *C.cCollectorT
}

func NewBridgeCollector(cCollector *C.cCollectorT) *bridgeCollector {
	return &bridgeCollector{
		cCollector: cCollector,
	}
}

func (bc *bridgeCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	return nil
}

func (bc *bridgeCollector) Collect(ctx plugin.CollectContext) error {
	C.callCollect(bc.cCollector.collectCallback)
	return nil
}

func (bc *bridgeCollector) Load(ctx plugin.Context) error {
	return nil
}

func (bc *bridgeCollector) Unload(ctx plugin.Context) error {
	return nil
}

// export StartCollector
func StartCollector(cCollector *C.cCollectorT, name *C.char, version *C.char) {
	bCollector := NewBridgeCollector(cCollector)
	runner.StartCollector(bCollector, C.GoString(name), C.GoString(version)) // todo: should release?
}

func main() {}
