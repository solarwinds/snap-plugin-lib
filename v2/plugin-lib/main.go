package main

import (
	"fmt"
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

/*****************************************************************************/
// helpers

func contextObject(ctxId *C.char) *proxy.PluginContext {
	id := C.GoString(ctxId)
	ctx, ok := contextMap.Load(id)
	if !ok {
		panic(fmt.Sprintf("can't aquire context object with id %s", id))
	}

	ctxObj, okType := ctx.(*proxy.PluginContext)
	if !okType {
		panic("Invalid concrete type of context object")
	}

	return ctxObj
}

func intToBool(v int) bool {
	return v != 0
}

/*****************************************************************************/
// Collect related functions

//export ctx_add_metric_int
func ctx_add_metric_int(ctxId *C.char, ns *C.char, v int) {
	contextObject(ctxId).AddMetric(C.GoString(ns), v)
}

/*****************************************************************************/
// DefinePlugin related functions

//export def_define_metric
func def_define_metric(namespace *C.char, unit *C.char, isDefault int, description *C.char) {
	pluginDef.DefineMetric(C.GoString(namespace), C.GoString(unit), intToBool(isDefault), C.GoString(description))
}

//export def_define_group
func def_define_group(name *C.char, description *C.char) {
	pluginDef.DefineGroup(C.GoString(name), C.GoString(description))
}

//export def_example_config
func def_example_config(cfg *C.char) {
	_ = pluginDef.DefineExampleConfig(C.GoString(cfg))
}

//export def_define_tasks_per_instance_limit
func def_define_tasks_per_instance_limit(limit int) {
	_ = pluginDef.DefineTasksPerInstanceLimit(limit)
}

//export def_define_instances_limit
func def_define_instances_limit(limit int) {
	_ = pluginDef.DefineInstancesLimit(limit)
}

/*****************************************************************************/

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

/*****************************************************************************/

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

/*****************************************************************************/

func main() {}
