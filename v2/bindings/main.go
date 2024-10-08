/*
 Copyright (c) 2024 SolarWinds Worldwide, LLC

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

//go:generate goversioninfo -64
package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <memory.h>

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
	TYPE_INT32,
	TYPE_UINT32,
	TYPE_FLOAT,
	TYPE_DOUBLE,
	TYPE_BOOL,
	TYPE_CSTRING,
	TYPE_INT16,
	TYPE_UINT16
};

typedef struct {
	union  {
		long long v_int64;
		unsigned long long v_uint64;
		int v_int32;
		unsigned int v_uint32;
		float v_float;
		double v_double;
		int v_bool;
		char * v_cstring;
		short int v_int16;
		unsigned short int v_uint16;
	} value;
	int vtype; // value_type_t;
} value_t;

static inline value_t * alloc_value_t(enum value_type_t t) {
	value_t * val_ptr = malloc(sizeof(value_t));
	val_ptr->vtype = t;
	return val_ptr;
}

static inline void free_value_t(value_t * v) {
	if (v->vtype == TYPE_CSTRING) {
		free(v->value.v_cstring);
	}
	free(v);
}

static inline long long value_t_long_long(value_t * v) { return v->value.v_int64; }
static inline unsigned long long value_t_ulong_long(value_t * v) { return v->value.v_uint64; }
static inline int value_t_int(value_t * v) { return v->value.v_int32; }
static inline unsigned int value_t_uint(value_t * v) { return v->value.v_uint32; }
static inline float value_t_float(value_t * v) { return v->value.v_float; }
static inline double value_t_double(value_t * v) { return v->value.v_double; }
static inline int value_t_bool(value_t * v) { return v->value.v_bool; }
static inline char * value_t_cstring(value_t * v) { return v->value.v_cstring; }
static inline short int value_t_shortint(value_t * v) { return v->value.v_int16; }
static inline short int value_t_ushortint(value_t * v) { return v->value.v_uint16; }

static inline void set_value_t_long_long(value_t * v, long long v_int64) { v->value.v_int64 = v_int64; }
static inline void set_value_t_ulong_long(value_t * v, unsigned long long v_uint64) { v->value.v_uint64 = v_uint64; }
static inline void set_value_t_int(value_t * v, int v_int32) { v->value.v_int32 = v_int32; }
static inline void set_value_t_uint(value_t * v, unsigned int v_uint32) { v->value.v_uint32 = v_uint32; }
static inline void set_value_t_float(value_t * v, float v_float) { v->value.v_float = v_float; }
static inline void set_value_t_double(value_t * v, double v_double) { v->value.v_double = v_double; }
static inline void set_value_t_bool(value_t * v, int v_bool) { v->value.v_bool = v_bool; }
static inline void set_value_t_cstring(value_t * v, char * v_cstring) { v->value.v_cstring = v_cstring; }
static inline void set_value_t_shortint(value_t * v, short int v_int16) { v->value.v_int16 = v_int16; }
static inline void set_value_t_ushortint(value_t * v, unsigned short int v_uint16) { v->value.v_uint16 = v_uint16; }

typedef struct {
	char * key;
	char * value;
} map_element_t;

static inline map_element_t * alloc_map_element_t_array(int size) {
	map_element_t * map_arr = malloc(sizeof(map_element_t) * size);
	return map_arr;
}

static inline void free_map_element_t_array(map_element_t* m, int size) {
	int i;
	for(i = 0; i < size; i++) {
		free(m[i].key);
		free(m[i].value);
	}
	free(m);
}

static inline void set_tag_values(map_element_t * tag_arr, int index, char * key, char * value) {
	tag_arr[index].key = key;
	tag_arr[index].value = value;
}

typedef struct {
	map_element_t * elements;
	int length;
} map_t;

static inline map_t * alloc_map_t() {
	map_t * map = malloc(sizeof(map_t));
	return map;
}

static inline void free_map_t(map_t * m) {
	free_map_element_t_array(m->elements, m->length);
	free(m);
}

static inline void set_map_elements(map_t * map_ptr, map_element_t * elements) {
	map_ptr->elements = elements;
}

static inline char * get_map_key(map_t * map, int index) { return map->elements[index].key; }
static inline char * get_map_value(map_t * map, int index) { return map->elements[index].value; }
static inline int get_map_length(map_t * map) { return map->length; }

static inline void set_map_lenght(map_t * map, int length) { map->length = length; }

typedef struct {
	char * msg;
} error_t;

static inline error_t * alloc_error_msg(char * msg) {
	error_t * errMsg = malloc(sizeof(error_t));
	errMsg->msg = msg;
	return errMsg;
}

static inline void free_error_msg(error_t * err) {
	if (err == NULL) return;

	if (err->msg != NULL) {
		free(err->msg);
		err->msg = NULL;
	}

	free(err);
}

typedef struct {
	int sec;
	int nsec;
} time_with_ns_t;

static inline time_with_ns_t* alloc_time_with_ns_t() {
	return malloc(sizeof(time_with_ns_t));
}

static inline void free_time_with_ns_t(time_with_ns_t* t) {
	free(t);
}

static inline void set_time_with_ns_t(time_with_ns_t* time_ptr, int sec, int nsec) {
	time_ptr->sec = sec;
	time_ptr->nsec = nsec;
}

enum metric_type_t {
	METRIC_TYPE_UNKNOWN,
	METRIC_TYPE_GAUGE,
	METRIC_TYPE_SUM,
	METRIC_TYPE_SUMMARY,
	METRIC_TYPE_HISTOGRAM
};

typedef struct {
	map_t * tags_to_add;
	map_t * tags_to_remove;
	time_with_ns_t * timestamp;
	char * description;
	char * unit;
	int metric_type;
} modifiers_t;


static inline char** alloc_str_array(int size) {
	return malloc(sizeof(char*) * size);
}

static inline void free_str_array(char **arr) {
	if (arr == NULL) return;

	char * arrEl = *arr;
	for (;;) {
		if (arrEl == NULL) {
			break;
		}

		free(arrEl);
		arrEl++;
	}

	free(arr);
}

static inline void set_str_array_element(char **str_array, int index, char *element) {
	str_array[index] = element;
}

typedef struct {
	char * el_name;
	char * value;
	char * description;
} namespace_element_t;


static inline namespace_element_t * alloc_namespace_elem_arr(int size) {
	namespace_element_t * ne_arr = malloc(sizeof(namespace_element_t) * size);
	return ne_arr;
}

static inline void free_namespace_elem_arr(namespace_element_t * nm_array, int size) {
	int i;
	for(i=0; i < size; i++){
		free(nm_array[i].el_name);
		free(nm_array[i].value);
		free(nm_array[i].description);
	}
	free(nm_array);
}

static inline void set_namespace_element(namespace_element_t * ne_arr, int index, char * el_name, char * value, char * description) {
	ne_arr[index].el_name = el_name;
	ne_arr[index].value = value;
	ne_arr[index].description = description;
}


typedef struct {
	namespace_element_t * nm_elements;
	int nm_length;
	char * nm_string;
} namespace_t;


static inline namespace_t * alloc_namespace_t() {
	namespace_t* nm_ptr = malloc(sizeof(namespace_t));
	return nm_ptr;
}

static inline void free_namespace_t(namespace_t * namespace_ptr) {
	free_namespace_elem_arr(namespace_ptr->nm_elements, namespace_ptr->nm_length);
	free(namespace_ptr);
}
static inline void set_namespace_fields(namespace_t * nm_ptr, namespace_element_t * ne_arr, int nm_length, char * nm_string) {
	nm_ptr->nm_elements = ne_arr;
	nm_ptr->nm_length = nm_length;
	nm_ptr->nm_string= nm_string;
}

typedef struct {
	namespace_t * mt_namespace;
	char * mt_description;
	value_t *mt_value;
	time_with_ns_t * timestamp;
	map_t * tags; // free
} metric_t;


static inline metric_t** alloc_metric_pointer_array(int size) {
	metric_t ** arrPtr = malloc(sizeof(metric_t*) * size);
	int i;
	for(i=0; i < size; i++) {
		arrPtr[i] = malloc(sizeof(metric_t));
	}
	return arrPtr;
}

static inline void set_metric_pointer_array_element(metric_t** mt_array, int index, metric_t* element) {
	mt_array[index] = element;
}

static inline void set_metric_values(metric_t** mt_array, int index, namespace_t* mt_namespace, char* desc, value_t* val, time_with_ns_t* timestamp, map_t* tags) {
	mt_array[index]->mt_namespace = mt_namespace;
	mt_array[index]->mt_description = desc;
	mt_array[index]->mt_value = val;
	mt_array[index]->timestamp = timestamp;
	mt_array[index]->tags = tags;
}

static inline void free_metric_arr(metric_t** mt_array, int size) {
	if (mt_array == NULL) return;
	int i;
	for (i=0; i< size; i++) {
		if (mt_array[i] != NULL ) {
			free_namespace_t(mt_array[i]->mt_namespace);
			free_value_t(mt_array[i]->mt_value);
			free_time_with_ns_t(mt_array[i]->timestamp);
			free_map_t(mt_array[i]->tags);
			free(mt_array[i]);
	   }
   }
   free(mt_array);
}
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
	collectorProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/collector/proxy"
	commonProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/proxy"
	publisherProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/publisher/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/runner"
)

var contextMap = sync.Map{}
var pluginDef plugin.Definition
var collectorDef plugin.CollectorDefinition

/*****************************************************************************/
// ctx helpers
func commonContextObject(ctxId *C.char) *commonProxy.Context {
	id := C.GoString(ctxId)
	ctx, ok := contextMap.Load(id)
	if !ok {
		panic(fmt.Sprintf("Can't aquire context object with id %v", id))
	}
	switch v := ctx.(type) {
	case *collectorProxy.PluginContext:
		ctxObj, okType := ctx.(*collectorProxy.PluginContext)
		if !okType {
			panic("Invalid concrete type of context object")
		}
		return ctxObj.Context
	case *publisherProxy.PluginContext:
		ctxObj, okType := ctx.(*publisherProxy.PluginContext)
		if !okType {
			panic("Invalid concrete type of context object")
		}
		return ctxObj.Context
	default:
		panic(fmt.Sprintf("%s not supported plugin context type", v))
	}

}

func publContextObject(ctxId *C.char) *publisherProxy.PluginContext {
	id := C.GoString(ctxId)
	ctx, ok := contextMap.Load(id)
	if !ok {
		panic(fmt.Sprintf("can't aquire context object with id %v", id))
	}

	ctxObj, okType := ctx.(*publisherProxy.PluginContext)
	if !okType {
		panic("Invalid concrete type of context object")
	}
	return ctxObj
}

func collContextObject(ctxId *C.char) *collectorProxy.PluginContext {
	id := C.GoString(ctxId)
	ctx, ok := contextMap.Load(id)
	if !ok {
		panic(fmt.Sprintf("can't aquire context object with id %v", id))
	}

	ctxObj, okType := ctx.(*collectorProxy.PluginContext)
	if !okType {
		panic("Invalid concrete type of context object")
	}
	return ctxObj
}

/*****************************************************************************/
// mapping helpers
func intToBool(v int) bool {
	return v != 0
}

func boolToInt(v bool) int {
	if !v {
		return 0
	}

	return 1
}

func toCmap_t(gomap map[string]string) *C.map_t {
	cMapPtr := C.alloc_map_t()
	map_len := len(gomap)
	C.set_map_lenght(cMapPtr, C.int(map_len))
	tagArrPtr := C.alloc_map_element_t_array(C.int(map_len))
	i := 0
	for key, val := range gomap {
		C.set_tag_values(tagArrPtr, C.int(i), (*C.char)(C.CString(key)), (*C.char)(C.CString(val)))
		i++
	}
	C.set_map_elements(cMapPtr, tagArrPtr)
	return cMapPtr
}

func toGoMap(m *C.map_t) map[string]string {
	tagsMap := map[string]string{}
	for i := 0; i < int(C.get_map_length(m)); i++ {
		k := C.GoString(C.get_map_key(m, C.int(i)))
		v := C.GoString(C.get_map_value(m, C.int(i)))
		tagsMap[k] = v
	}
	return tagsMap
}

func toCError(err error) *C.error_t {
	var errMsg *C.char
	if err != nil {
		errMsg = (*C.char)(C.CString(err.Error()))
	}
	return C.alloc_error_msg((*C.char)(errMsg))
}

func time_to_ctimewithns(timestamp time.Time) *C.time_with_ns_t {
	time_struct_ptr := C.alloc_time_with_ns_t()
	sec := timestamp.Unix()
	nsec := timestamp.UnixNano() - (timestamp.Unix() * 1e9)
	C.set_time_with_ns_t(time_struct_ptr, C.int(sec), C.int(nsec))
	return time_struct_ptr
}

func toCStrArray(arr []string) **C.char {
	cStrArr := C.alloc_str_array(C.int(len(arr) + 1))
	for i, el := range arr {
		C.set_str_array_element(cStrArr, C.int(i), (*C.char)(C.CString(el)))
	}
	C.set_str_array_element(cStrArr, C.int(len(arr)), (*C.char)(C.NULL)) // last element of array is always None (NULL)

	return cStrArr
}

func toCNamespace_t(nm plugin.Namespace) *C.namespace_t {
	nm_ptr := C.alloc_namespace_t()
	ne_arr := C.alloc_namespace_elem_arr(C.int(nm.Len()))
	for i := 0; i < nm.Len(); i++ {
		el := nm.At(i)
		C.set_namespace_element(ne_arr, C.int(i), (*C.char)(C.CString(el.Name())), (*C.char)(C.CString(el.Value())), (*C.char)(C.CString(el.Description())))
	}
	C.set_namespace_fields(nm_ptr, ne_arr, C.int(nm.Len()), (*C.char)(C.CString(nm.String())))
	return nm_ptr
}

func toCvalue_t(v interface{}) *C.value_t {
	switch n := v.(type) {
	case int64:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_INT64)
		C.set_value_t_long_long(cvalue_t_ptr, C.longlong(n))
		return cvalue_t_ptr
	case uint64:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_UINT64)
		C.set_value_t_ulong_long(cvalue_t_ptr, C.ulonglong(n))
		return cvalue_t_ptr
	case int32:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_INT32)
		C.set_value_t_int(cvalue_t_ptr, C.int(n))
		return cvalue_t_ptr
	case uint32:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_UINT32)
		C.set_value_t_uint(cvalue_t_ptr, C.uint(n))
		return cvalue_t_ptr
	case float32:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_FLOAT)
		C.set_value_t_float(cvalue_t_ptr, C.float(n))
		return cvalue_t_ptr
	case float64:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_DOUBLE)
		C.set_value_t_double(cvalue_t_ptr, C.double(n))
		return cvalue_t_ptr
	case bool:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_BOOL)
		boolint := 0
		if n {
			boolint = 1
		}
		C.set_value_t_bool(cvalue_t_ptr, C.int(boolint))
		return cvalue_t_ptr
	case string:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_CSTRING)
		C.set_value_t_cstring(cvalue_t_ptr, C.CString(n))
		return cvalue_t_ptr
	case int16:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_INT16)
		C.set_value_t_shortint(cvalue_t_ptr, C.short(n))
		return cvalue_t_ptr
	case uint16:
		cvalue_t_ptr := C.alloc_value_t(C.TYPE_UINT16)
		C.set_value_t_ushortint(cvalue_t_ptr, C.ushort(n))
		return cvalue_t_ptr
	default:
		panic(fmt.Sprintf("Not supported metric type %T", v))
	}
}

func toGoValue(v *C.value_t) interface{} {
	switch (*v).vtype {
	case C.TYPE_INT64:
		return int(C.value_t_long_long(v))
	case C.TYPE_UINT64:
		return uint(C.value_t_ulong_long(v))
	case C.TYPE_INT32:
		return int32(C.value_t_int(v))
	case C.TYPE_UINT32:
		return uint32(C.value_t_uint(v))
	case C.TYPE_FLOAT:
		return float32(C.value_t_float(v))
	case C.TYPE_DOUBLE:
		return float64(C.value_t_double(v))
	case C.TYPE_BOOL:
		return intToBool(int(C.value_t_bool(v)))
	case C.TYPE_CSTRING:
		return C.GoString(C.value_t_cstring(v))
	case C.TYPE_INT16:
		return int16(C.value_t_shortint(v))
	case C.TYPE_UINT16:
		return uint16(C.value_t_ushortint(v))
	}

	panic(fmt.Sprintf("Invalid type %v", (*v).vtype))
}

func toGoModifiers(modifiers *C.modifiers_t) []plugin.MetricModifier {
	var appliedModifiers []plugin.MetricModifier

	if modifiers == nil {
		return appliedModifiers
	}

	if modifiers.tags_to_add != nil {
		appliedModifiers = append(appliedModifiers, plugin.MetricTags(toGoMap(modifiers.tags_to_add)))
	}

	if modifiers.tags_to_remove != nil {
		appliedModifiers = append(appliedModifiers, plugin.MetricTags(toGoMap(modifiers.tags_to_remove)))
	}

	if modifiers.timestamp != nil {
		appliedModifiers = append(appliedModifiers,
			plugin.MetricTimestamp(time.Unix(int64(modifiers.timestamp.sec), int64(modifiers.timestamp.nsec))))
	}

	if modifiers.description != nil {
		appliedModifiers = append(appliedModifiers, plugin.MetricDescription(C.GoString(modifiers.description)))
	}

	if modifiers.unit != nil {
		appliedModifiers = append(appliedModifiers, plugin.MetricUnit(C.GoString(modifiers.unit)))
	}

	if modifiers.metric_type != C.METRIC_TYPE_UNKNOWN {
		appliedModifiers = append(appliedModifiers, toGoMetricTypeModifier(modifiers.metric_type))
	}

	return appliedModifiers
}

func toGoMetricTypeModifier(metricType C.int) plugin.MetricModifier {
	switch metricType {
	case C.METRIC_TYPE_GAUGE:
		return plugin.MetricTypeGauge()
	case C.METRIC_TYPE_SUM:
		return plugin.MetricTypeSum()
	case C.METRIC_TYPE_SUMMARY:
		return plugin.MetricTypeSummary()
	case C.METRIC_TYPE_HISTOGRAM:
		return plugin.MetricTypeHistogram()
	default:
		panic(fmt.Sprintf("Invalid metric type %v", metricType))
	}
}

/*****************************************************************************/
// Deallocate API

//export dealloc_charp
func dealloc_charp(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p)) // #nosec G103
	}
}

//export dealloc_str_array
func dealloc_str_array(p **C.char) {
	C.free_str_array(p)
}

//export dealloc_error
func dealloc_error(p *C.error_t) {
	C.free_error_msg(p)
}

//export dealloc_metric_array
func dealloc_metric_array(p **C.metric_t, size int) {
	C.free_metric_arr(p, C.int(size))
}

/*****************************************************************************/
// C API - Collect related functions

// ctx.Store() and ctx.Load() have to be implemented and managed on the
// native language side due to garbage collection functionality (we can't
// pass Python/C# address to Go side, since it may become invalid).

// ctx.Done() returns channel which is not a simple type present in other
// language. Only ctx.IsDone() may be used

//export ctx_add_metric
func ctx_add_metric(ctxID *C.char, ns *C.char, v *C.value_t, modifiers *C.modifiers_t) *C.error_t {
	err := collContextObject(ctxID).AddMetric(C.GoString(ns), toGoValue(v), toGoModifiers(modifiers)...)
	return toCError(err)
}

//export ctx_always_apply
func ctx_always_apply(ctxID *C.char, ns *C.char, modifiers *C.modifiers_t) *C.error_t {
	_, err := collContextObject(ctxID).AlwaysApply(C.GoString(ns), toGoModifiers(modifiers)...)
	return toCError(err)
}

//export ctx_dismiss_all_modifiers
func ctx_dismiss_all_modifiers(ctxID *C.char) {
	collContextObject(ctxID).DismissAllModifiers()
}

//export ctx_should_process
func ctx_should_process(ctxID *C.char, ns *C.char) int {
	return boolToInt(collContextObject(ctxID).ShouldProcess(C.GoString(ns)))
}

//export ctx_requested_metrics
func ctx_requested_metrics(ctxID *C.char) **C.char {
	return toCStrArray(collContextObject(ctxID).RequestedMetrics())
}

//export ctx_count
func ctx_count(ctxID *C.char) int {
	return publContextObject(ctxID).Count()
}

//export ctx_list_all_metrics
func ctx_list_all_metrics(ctxID *C.char) **C.metric_t {
	mts := publContextObject(ctxID).ListAllMetrics()
	mtPtArr := C.alloc_metric_pointer_array(C.int(len(mts)))

	for i, el := range mts {
		mtNamespace := toCNamespace_t(el.Namespace())
		mtDesc := (*C.char)(C.CString(el.Description()))
		mtValue := toCvalue_t(el.Value())
		mtTimestamp := time_to_ctimewithns(el.Timestamp())
		mtTags := toCmap_t(el.Tags())
		C.set_metric_values(mtPtArr, C.int(i), mtNamespace, mtDesc, mtValue, mtTimestamp, mtTags)
	}
	return mtPtArr
}

///////////////////////////////////////////////////////////////////////////////

//export ctx_config_value
func ctx_config_value(ctxID *C.char, key *C.char) *C.char {
	v, ok := commonContextObject(ctxID).ConfigValue(C.GoString(key))

	if !ok {
		return (*C.char)(C.NULL)
	}

	return C.CString(v)
}

//export ctx_config_keys
func ctx_config_keys(ctxID *C.char) **C.char {
	return toCStrArray(commonContextObject(ctxID).ConfigKeys())
}

//export ctx_raw_config
func ctx_raw_config(ctxID *C.char) *C.char {
	rc := string(commonContextObject(ctxID).RawConfig())
	return C.CString(rc)
}

//export ctx_add_warning
func ctx_add_warning(ctxID *C.char, message *C.char) {
	commonContextObject(ctxID).AddWarning(C.GoString(message))
}

//export ctx_is_done
func ctx_is_done(ctxID *C.char) int {
	return boolToInt(commonContextObject(ctxID).IsDone())
}

//export ctx_log
func ctx_log(ctxID *C.char, level C.int, message *C.char, fields *C.map_t) {
	logF := commonContextObject(ctxID).Logger()

	if fields != nil {
		for i := 0; i < int(C.get_map_length(fields)); i++ {
			k := C.get_map_key(fields, C.int(i))
			v := C.get_map_value(fields, C.int(i))
			logF = logF.WithField(C.GoString(k), C.GoString(v))
		}
	}

	if logObj, ok := logF.(*logrus.Entry); ok {
		logObj.Log(logrus.Level(int(level)), C.GoString(message))
	}
}

/*****************************************************************************/
// C API - DefinePlugin related functions

//export define_metric
func define_metric(namespace *C.char, unit *C.char, isDefault int, description *C.char) {
	collectorDef.DefineMetric(C.GoString(namespace), C.GoString(unit), intToBool(isDefault), C.GoString(description))
}

//export define_group
func define_group(name *C.char, description *C.char) {
	collectorDef.DefineGroup(C.GoString(name), C.GoString(description))
}

//export define_example_config
func define_example_config(cfg *C.char) *C.error_t {
	err := collectorDef.DefineExampleConfig(C.GoString(cfg))
	return toCError(err)
}

///////////////////////////////////////////////////////////////////////////////

//export define_tasks_per_instance_limit
func define_tasks_per_instance_limit(limit int) {
	_ = pluginDef.DefineTasksPerInstanceLimit(limit)
}

//export define_instances_limit
func define_instances_limit(limit int) {
	_ = pluginDef.DefineInstancesLimit(limit)
}

/*****************************************************************************/
// C API - runner related functions

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

//export start_streaming_collector
func start_streaming_collector(collectCallback *C.callback_t, loadCallback *C.callback_t, unloadCallback *C.callback_t, defineCallback *C.define_callback_t, name *C.char, version *C.char) {
	bCollector := &bridgeCollector{
		collectCallback: collectCallback,
		loadCallback:    loadCallback,
		unloadCallback:  unloadCallback,
		defineCallback:  defineCallback,
	}
	runner.StartStreamingCollector(bCollector, C.GoString(name), C.GoString(version))
}

/***************************************************************************/
//export start_publisher
func start_publisher(publishCallback *C.callback_t, loadCallback *C.callback_t, unloadCallback *C.callback_t, defineCallback *C.define_callback_t, name *C.char, version *C.char) {

	bPublisher := &bridgePublisher{
		publishCallback: publishCallback,
		loadCallback:    loadCallback,
		unloadCallback:  unloadCallback,
		defineCallback:  defineCallback,
	}
	runner.StartPublisher(bPublisher, C.GoString(name), C.GoString(version))
}

/****************************************************************************/

type bridgePublisher struct {
	publishCallback *C.callback_t
	loadCallback    *C.callback_t
	unloadCallback  *C.callback_t
	defineCallback  *C.define_callback_t
}

func (bp *bridgePublisher) Publish(ctx plugin.PublishContext) error {
	return bp.callC(ctx, bp.publishCallback)
}

func (bp *bridgePublisher) PluginDefinition(def plugin.PublisherDefinition) error {
	pluginDef = def
	C.call_c_define_callback(bp.defineCallback)

	return nil
}

func (bp *bridgePublisher) Load(ctx plugin.Context) error {
	return bp.callC(ctx, bp.loadCallback)
}

func (bp *bridgePublisher) Unload(ctx plugin.Context) error {
	return bp.callC(ctx, bp.unloadCallback)
}

func (bp *bridgePublisher) callC(ctx plugin.Context, callback *C.callback_t) error {
	ctxAsType := ctx.(*publisherProxy.PluginContext)
	taskID := ctxAsType.TaskID()
	contextMap.Store(taskID, ctxAsType)
	defer contextMap.Delete(taskID)

	taskIDasPtr := C.CString(taskID)
	C.call_c_callback(callback, taskIDasPtr)
	dealloc_charp(taskIDasPtr)

	return nil
}

type bridgeCollector struct {
	collectCallback *C.callback_t
	loadCallback    *C.callback_t
	unloadCallback  *C.callback_t
	defineCallback  *C.define_callback_t
}

func (bc *bridgeCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	pluginDef = def
	collectorDef = def
	C.call_c_define_callback(bc.defineCallback)

	return nil
}

func (bc *bridgeCollector) Collect(ctx plugin.CollectContext) error {
	return bc.callC(ctx, bc.collectCallback)
}

func (bc *bridgeCollector) StreamingCollect(ctx plugin.CollectContext) error {
	return bc.callC(ctx, bc.collectCallback)
}

func (bc *bridgeCollector) Load(ctx plugin.Context) error {
	return bc.callC(ctx, bc.loadCallback)
}

func (bc *bridgeCollector) Unload(ctx plugin.Context) error {
	return bc.callC(ctx, bc.unloadCallback)
}

func (bc *bridgeCollector) callC(ctx plugin.Context, callback *C.callback_t) error {
	ctxAsType := ctx.(*collectorProxy.PluginContext)
	taskID := ctxAsType.TaskID()

	contextMap.Store(taskID, ctxAsType)
	defer contextMap.Delete(taskID)

	taskIDasPtr := C.CString(taskID)
	C.call_c_callback(callback, taskIDasPtr)

	dealloc_charp(taskIDasPtr)

	return nil
}

/*****************************************************************************/

func main() {}
