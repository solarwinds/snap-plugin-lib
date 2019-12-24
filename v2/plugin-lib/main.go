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
enum value_type_t {
	TYPE_INVALID,
	TYPE_INT64,
	TYPE_UINT64,
	TYPE_DOUBLE,
	TYPE_BOOL,
};

typedef struct {
	union  {
		long long v_int64;
		unsigned long long v_uint64;
		double v_double;
		int v_bool;
	} value;
	int vtype; // value_type_t;
} value_t;

static inline long long value_t_long_long(value_t * v) { return v->value.v_int64; }
static inline unsigned long long value_t_ulong_long(value_t * v) { return v->value.v_uint64; }
static inline double value_t_double(value_t * v) { return v->value.v_double; }
static inline int value_t_bool(value_t * v) { return v->value.v_bool; }

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

static inline void free_error_msg(error_t * err) {
	if (err == NULL) {
		return;
	}

	if (err->msg != NULL) {
		free(err->msg);
	}

	free(err);
}

static inline void free_memory(void * p) {
	free(p);
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

func ctagsToMap(tags *C.tag_t, tagsCount int) map[string]string {
	tagsMap := map[string]string{}
	for i := 0; i < tagsCount; i++ {
		k := C.GoString(C.tag_key(tags, C.int(i)))
		v := C.GoString(C.tag_value(tags, C.int(i)))
		tagsMap[k] = v
	}
	return tagsMap
}

func toCError(err error) *C.error_t {
	var errMsg *C.char
	if err != nil {
		errMsg = (* C.char)(C.CString(err.Error()))
	}
	return C.alloc_error_msg((* C.char)(errMsg))
}

func toGoValue(v *C.value_t) interface{} {
	switch (*v).vtype {
	case C.TYPE_INT64:
		return int(C.value_t_long_long(v))
	case C.TYPE_UINT64:
		return uint(C.value_t_ulong_long(v))
	case C.TYPE_DOUBLE:
		return float64(C.value_t_double(v))
	case C.TYPE_BOOL:
		return intToBool(int(C.value_t_bool(v)))
	}

	panic("Invalid type")
}

/*****************************************************************************/
// Collect related functions

//export ctx_add_metric
func ctx_add_metric(ctx_id *C.char, ns *C.char, v int) *C.error_t {
	err := contextObject(ctx_id).AddMetric(C.GoString(ns), v)
	return toCError(err)
}

//export ctx_add_metric_ex
func ctx_add_metric_ex(ctx_id *C.char, ns *C.char, v *C.value_t) *C.error_t {
	err := contextObject(ctx_id).AddMetric(C.GoString(ns), toGoValue(v))
	return toCError(err)
}

//export ctx_add_metric_with_tags
func ctx_add_metric_with_tags(ctxId *C.char, ns *C.char, v int, tags *C.tag_t, tagsCount int) *C.error_t {
	err := contextObject(ctxId).AddMetricWithTags(C.GoString(ns), v, ctagsToMap(tags, tagsCount))
	return toCError(err)
}

//export ctx_apply_tags_by_path
func ctx_apply_tags_by_path(ctxId *C.char, ns *C.char, tags *C.tag_t, tagsCount int) *C.error_t {
	err := contextObject(ctxId).ApplyTagsByPath(C.GoString(ns), ctagsToMap(tags, tagsCount))
	return toCError(err)
}

//export ctx_apply_tags_by_regexp
func ctx_apply_tags_by_regexp(ctxId *C.char, ns *C.char, tags *C.tag_t, tagsCount int) *C.error_t {
	err := contextObject(ctxId).ApplyTagsByRegExp(C.GoString(ns), ctagsToMap(tags, tagsCount))
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
	rc := string(contextObject(ctxId).RawConfig())
	fmt.Printf("rc=%#v\n", rc)
	return C.CString(rc)
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

//export define_global_tags
func define_global_tags(namespace *C.char, tags *C.tag_t, tagsCount int) {
	pluginDef.DefineGlobalTags(C.GoString(namespace), ctagsToMap(tags, tagsCount))
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
