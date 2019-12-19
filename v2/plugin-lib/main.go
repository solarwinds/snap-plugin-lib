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

typedef void (callbackT)(char *); // used for Collect, Load and Unload
typedef void (defineCallbackT)(); // used for DefineCallback

// called from Go code
static inline void CCallback(callbackT callback, char * ctxId) { callback(ctxId); }
static inline void CDefineCallback(defineCallbackT callback) { callback(); }

typedef struct {
	char * key;
	char * value;
} tag;

static inline char * tag_key(tag tags, int index) { return tags.key; }
static inline char * tag_value(tag tags, int index) { return tags.value; }

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

func ctagsToMap(tags *C.tag, tagsCount int) map[string]string {
	tagsMap := map[string]string{}
	return tagsMap
	//
	//for i := 0; i < tagsCount; i++ {
	//	k := C.GoString(C.tag_key(tags, C.int(i)))
	//	v := C.GoString(C.tag_value(tags, C.int(i)))
	//	fmt.Printf("**********TAGS&*&&&&&&&&&& k=%#v\n", k)
	//	fmt.Printf("**********TAGS&*&&&&&&&&&& v=%#v\n", k)
	//	tagsMap[k] = v
	//}
	//return tagsMap
}

func ctagsToMap1(tags C.tag, tagsCount int) map[string]string {
	tagsMap := map[string]string{}
	for i := 0; i < 1; i++ {
		k := C.GoString(C.tag_key(tags, C.int(i)))
		v := C.GoString(C.tag_value(tags, C.int(i)))
		fmt.Printf("**********TAGS&*&&&&&&&&&& k=%#v\n", k)
		fmt.Printf("**********TAGS&*&&&&&&&&&& v=%#v\n", k)
		tagsMap[k] = v
	}
	return tagsMap
}

/*****************************************************************************/
// Collect related functions

//export ctx_add_metric
func ctx_add_metric(ctxId *C.char, ns *C.char, v int) {
	contextObject(ctxId).AddMetric(C.GoString(ns), v)
}

//export ctx_add_metric_with_tags
func ctx_add_metric_with_tags(ctxId *C.char, ns *C.char, v int, tags C.tag, tagsCount int) {
	contextObject(ctxId).AddMetricWithTags(C.GoString(ns), v, ctagsToMap1(tags, tagsCount))
}

//export ctx_apply_tags_by_path
func ctx_apply_tags_by_path(ctxId *C.char, ns *C.char, tags *C.tag, tagsCount int) {
	contextObject(ctxId).ApplyTagsByPath(C.GoString(ns), ctagsToMap(tags, tagsCount))
}

//export ctx_apply_tags_by_regexp
func ctx_apply_tags_by_regexp(ctxId *C.char, ns *C.char, tags *C.tag, tagsCount int) {
	contextObject(ctxId).ApplyTagsByRegExp(C.GoString(ns), ctagsToMap(tags, tagsCount))
}

//export ctx_should_process
func ctx_should_process(ctxId *C.char, ns *C.char) int {
	return boolToInt(contextObject(ctxId).ShouldProcess(C.GoString(ns)))
}

//export ctx_config
func ctx_config(ctxId *C.char, key *C.char) *C.char {
	v, _ := contextObject(ctxId).Config(C.GoString(key))
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
func define_example_config(cfg *C.char) {
	_ = pluginDef.DefineExampleConfig(C.GoString(cfg))
}

//export define_tasks_per_instance_limit
func define_tasks_per_instance_limit(limit int) {
	_ = pluginDef.DefineTasksPerInstanceLimit(limit)
}

//export define_instances_limit
func define_instances_limit(limit int) {
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
