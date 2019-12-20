package main

import (
	"fmt"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"reflect"
	"sync"
	"unsafe"
)

/*
#include <stdlib.h>

// c types for callbacks
typedef void (callback_t)(char *);  // used for Collect, Load and Unload
typedef void (define_callback_t)(); // used for DefineCallback

// called from Go code
static inline void call_c_callback(callback_t callback, char * ctxId) { callback(ctxId); }
static inline void call_c_define_callback(define_callback_t callback) { callback(); }

// some helpers to manage C/Go memory/access interactions
typedef struct {
    char * key;
    char * value;
} tag_t;

static inline char * tag_key(tag_t * tags, int index) { return tags[index].key; }
static inline char * tag_value(tag_t * tags, int index) { return tags[index].value; }

typedef struct {
    char * msg;
} error_t;

static inline error_t * alloc_error_msg(char * msg) {
    error_t * errMsg = malloc(sizeof(error_t));
    errMsg->msg = msg;
    return errMsg;
}

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

func boolToInt(v bool) int {
	if v == false {
		return 0
	}

	return 1
}

func ctagsToMap(tag_ts *C.tag_t, tag_tsCount int) map[string]string {
	tag_tsMap := map[string]string{}
	for i := 0; i < tag_tsCount; i++ {
		k := C.GoString(C.tag_key(tag_ts, C.int(i)))
		v := C.GoString(C.tag_value(tag_ts, C.int(i)))
		tag_tsMap[k] = v
	}
	return tag_tsMap
}

func toCError(err error) *C.error_t {
	var errMsg *C.char
	if err != nil {
		errMsg = (* C.char)(C.CString(err.Error()))
	}
	return C.alloc_error_msg((* C.char)(errMsg))
}

/*****************************************************************************/
// Collect related functions

//export ctx_add_metric
func ctx_add_metric(ctxId *C.char, ns *C.char, v int) *C.error_t {
	err := contextObject(ctxId).AddMetric(C.GoString(ns), v)
	return toCError(err)
}

//export ctx_add_metric_with_tag_ts
func ctx_add_metric_with_tag_ts(ctxId *C.char, ns *C.char, v int, tag_ts *C.tag_t, tag_tsCount int) *C.error_t {
	err := contextObject(ctxId).AddMetricWithTags(C.GoString(ns), v, ctagsToMap(tag_ts, tag_tsCount))
	return toCError(err)
}

//export ctx_apply_tag_ts_by_path
func ctx_apply_tag_ts_by_path(ctxId *C.char, ns *C.char, tag_ts *C.tag_t, tag_tsCount int) *C.error_t {
	err := contextObject(ctxId).ApplyTagsByPath(C.GoString(ns), ctagsToMap(tag_ts, tag_tsCount))
	return toCError(err)
}

//export ctx_apply_tag_ts_by_regexp
func ctx_apply_tag_ts_by_regexp(ctxId *C.char, ns *C.char, tag_ts *C.tag_t, tag_tsCount int) *C.error_t {
	err := contextObject(ctxId).ApplyTagsByRegExp(C.GoString(ns), ctagsToMap(tag_ts, tag_tsCount))
	return toCError(err)
}

//export ctx_should_process
func ctx_should_process(ctxId *C.char, ns *C.char) int {
	return boolToInt(contextObject(ctxId).ShouldProcess(C.GoString(ns)))
}

//export ctx_config
func ctx_config(ctxId *C.char, key *C.char) *C.char {
	v, ok := contextObject(ctxId).Config(C.GoString(key))
	if !ok {
		return (* C.char)(C.NULL)
	}

	return C.CString(v)
}

// todo: ctx_config_keys

//export ctx_raw_config
func ctx_raw_config(ctxId *C.char) *C.char {
	return C.CString(string(contextObject(ctxId).RawConfig()))
}

//export ctx_store
func ctx_store(ctxId *C.char, key *C.char, obj unsafe.Pointer) {
	contextObject(ctxId).Store(C.GoString(key), obj)
}

//export ctx_load
func ctx_load(ctxId *C.char, key *C.char) unsafe.Pointer {
	v, _ := contextObject(ctxId).Load(C.GoString(key))
	return unsafe.Pointer(reflect.ValueOf(v).Pointer())
}

/*****************************************************************************/
// DefinePlugin related functions

//export define_metric
func define_metric(namespace *C.char, unit *C.char, isDefault int, description *C.char) {
	pluginDef.DefineMetric(C.GoString(namespace), C.GoString(unit), intToBool(isDefault), C.GoString(description))
}

//export define_group
func define_group(name *C.char, description *C.char) {
	pluginDef.DefineGroup(C.GoString(name), C.GoString(description))
}

//export define_example_config
func define_example_config(cfg *C.char) *C.error_t {
	err := pluginDef.DefineExampleConfig(C.GoString(cfg))
	return toCError(err)
}

//export define_tasks_per_instance_limit
func define_tasks_per_instance_limit(limit int) {
	pluginDef.DefineTasksPerInstanceLimit(limit)
}

//export define_instances_limit
func define_instances_limit(limit int) {
	pluginDef.DefineInstancesLimit(limit)
}

/*****************************************************************************/

//export start_collector
func start_collector(collectCallback *C.callback_t, loadCallback *C.callback_t, unloadCallback *C.callback_t, defineCallback *C.define_callback_t, name *C.char, version *C.char) {
	bCollector := &bridgeCollector{
		collectCallback: collectCallback,
		loadCallback:    loadCallback,
		unloadCallback:  unloadCallback,
		defineCallback:  defineCallback,
	}
	runner.StartCollector(bCollector, C.GoString(name), C.GoString(version))
}

/*****************************************************************************/

type bridgeCollector struct {
	collectCallback *C.callback_t
	loadCallback    *C.callback_t
	unloadCallback  *C.callback_t
	defineCallback  *C.define_callback_t
}

func (bc *bridgeCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	pluginDef = def
	C.call_c_define_callback(bc.defineCallback)

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

func (bc *bridgeCollector) callC(ctx plugin.Context, callback *C.callback_t) error {
	ctxAsType := ctx.(*proxy.PluginContext)
	taskID := ctxAsType.TaskID()

	contextMap.Store(taskID, ctxAsType)
	defer contextMap.Delete(taskID)

	C.call_c_callback(callback, C.CString(taskID))
	return nil
}

/*****************************************************************************/

func main() {}
