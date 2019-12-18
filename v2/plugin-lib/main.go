package main

import (
	"fmt"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"unsafe"
)

/*
#include <stdlib.h>

typedef void (collectCallbackT)(void *);

typedef struct {
	collectCallbackT *collectCallback;
} cCollectorT;

// called from Go code
static inline void Collect(collectCallbackT collectCallback, void * ctx) { collectCallback(ctx); }

*/
import "C"

//export ctx_add_metric
func ctx_add_metric(ctx unsafe.Pointer, ns *C.char) {
	//ctxC := (* mock.Context)(ctx)
	//ctxC.AddMetric(C.GoString(ns), 10)
	fmt.Printf("*************I'm hehe***************")
}

type bridgeCollector struct {
	collectCallback * C.collectCallbackT
}

func (bc *bridgeCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	return nil
}

func (bc *bridgeCollector) Collect(ctx plugin.CollectContext) error {
	//ptr := unsafe.Pointer(reflect.ValueOf(ctx).Pointer())

	fmt.Printf("***GOADAMIK1\n")
	cptr := bc.collectCallback
	fmt.Printf("cptr=%#v\n", cptr)
	C.Collect(cptr, nil)
	fmt.Printf("***GOADAMIK2\n")

	return nil
}

func (bc *bridgeCollector) Load(ctx plugin.Context) error {
	return nil
}

func (bc *bridgeCollector) Unload(ctx plugin.Context) error {
	return nil
}

//export StartCollector
func StartCollector(callback * C.collectCallbackT, name *C.char, version *C.char) {
	bCollector := &bridgeCollector{collectCallback:callback}
	runner.StartCollector(bCollector, C.GoString(name), C.GoString(version)) // todo: should release?
}

func main() {}
