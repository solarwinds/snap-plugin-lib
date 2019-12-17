package main

import (
	"fmt"
	"github.com/librato/snap-plugin-lib-go/v2/mock"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"reflect"
	"unsafe"
)

/*
#include <stdlib.h>

typedef char * (collectCallbackT)(void *);

typedef struct {
	collectCallbackT *collectCallback;
} cCollectorT;


// called from Go code
static inline void Collect(collectCallbackT collectCallback, void * ctx) { collectCallback(ctx); }


*/
import "C"

//export ctx_add_metric
func ctx_add_metric(ctx unsafe.Pointer, ns * C.char ) {
	ctxC := (* mock.Context)(ctx)
	ctxC.AddMetric(C.GoString(ns), 10)
}

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
	ptr := unsafe.Pointer(reflect.ValueOf(ctx).Pointer())
	fmt.Printf("ptr=%#v\n", ptr)

	C.Collect(bc.cCollector.collectCallback, ptr)
	return nil
}

func (bc *bridgeCollector) Load(ctx plugin.Context) error {
	return nil
}

func (bc *bridgeCollector) Unload(ctx plugin.Context) error {
	return nil
}

//export StartCollector
func StartCollector(cCollector *C.cCollectorT, name *C.char, version *C.char) {
	bCollector := NewBridgeCollector(cCollector)
	runner.StartCollector(bCollector, C.GoString(name), C.GoString(version)) // todo: should release?
}

func main() {}
