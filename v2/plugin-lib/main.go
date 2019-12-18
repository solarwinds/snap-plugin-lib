package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"sync"
)

/*
#include <stdlib.h>

typedef void (collectCallbackT)(void *);

typedef struct {
	collectCallbackT *collectCallback;
} cCollectorT;

// called from Go code
static inline void Collect(collectCallbackT collectCallback, char * ctxId) { collectCallback(ctxId); }

*/
import "C"

var contextMap sync.Map = sync.Map{}

//export ctx_add_metric_int
func ctx_add_metric_int(ctxId *C.char, ns *C.char, v int) {
	id := C.GoString(ctxId)
	ctx, _ := contextMap.Load(id)
	ctx.(* proxy.PluginContext).AddMetric(C.GoString(ns), v)
}

type bridgeCollector struct {
	collectCallback *C.collectCallbackT
}

func (bc *bridgeCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	return nil
}

func (bc *bridgeCollector) Collect(ctx plugin.CollectContext) error {
	ctxAsType := ctx.(*proxy.PluginContext)
	taskID := ctxAsType.TaskID()
	contextMap.Store(taskID, ctxAsType)
	C.Collect(bc.collectCallback, C.CString(taskID))
	contextMap.Delete(taskID)

	return nil
}

func (bc *bridgeCollector) Load(ctx plugin.Context) error {
	return nil
}

func (bc *bridgeCollector) Unload(ctx plugin.Context) error {
	return nil
}

//export StartCollector
func StartCollector(callback *C.collectCallbackT, name *C.char, version *C.char) {
	bCollector := &bridgeCollector{collectCallback: callback}
	runner.StartCollector(bCollector, C.GoString(name), C.GoString(version)) // todo: should release?
}

func main() {}
