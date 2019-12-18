package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"sync"
)

/*
#include <stdlib.h>

typedef void (callbackT)(char *); // used for Collect, Load and Unload
typedef void (defineCallbackT)(); // used for DefineCallback

// called from Go code
static inline void CCallback(callbackT callback, char * ctxId) { callback(ctxId); }
static inline void CDefineCallback(defineCallbackT callback) { callback(); }

*/
import "C"

var contextMap = sync.Map{}
var pluginDef plugin.CollectorDefinition

//export ctx_add_metric_int
func ctx_add_metric_int(ctxId *C.char, ns *C.char, v int) {
	ctx, _ := contextMap.Load(C.GoString(ctxId))
	ctx.(*proxy.PluginContext).AddMetric(C.GoString(ns), v)
}

//export StartCollector
func StartCollector(collectCallback *C.callbackT, loadCallback *C.callbackT, unloadCallback *C.callbackT, defineCallback *C.defineCallbackT, name *C.char, version *C.char) {
	bCollector := &bridgeCollector{
		collectCallback: collectCallback,
		loadCallback:    loadCallback,
		unloadCallback:  unloadCallback,
		defineCallback:  defineCallback,
	}
	runner.StartCollector(bCollector, C.GoString(name), C.GoString(version)) // todo: should release?
}

type bridgeCollector struct {
	collectCallback *C.callbackT
	loadCallback    *C.callbackT
	unloadCallback  *C.callbackT
	defineCallback  *C.defineCallbackT
}

func (bc *bridgeCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	pluginDef = def
	C.CDefineCallback(bc.defineCallback)

	return nil
}

func (bc *bridgeCollector) Collect(ctx plugin.CollectContext) error {
	return bc.callC(ctx, bc.collectCallback)
}

func (bc *bridgeCollector) Load(ctx plugin.Context) error {
	return bc.callC(ctx, bc.loadCallback)
}

func (bc *bridgeCollector) Unload(ctx plugin.Context) error {
	return bc.callC(ctx, bc.unloadCallback)
}

func (bc *bridgeCollector) callC(ctx plugin.Context, callback *C.callbackT) error {
	ctxAsType := ctx.(*proxy.PluginContext)
	taskID := ctxAsType.TaskID()

	contextMap.Store(taskID, ctxAsType)
	defer contextMap.Delete(taskID)

	C.CCallback(callback, C.CString(taskID))
	return nil
}

func main() {}
