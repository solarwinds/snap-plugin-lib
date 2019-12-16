package main

import "C"


/*
#include <stdlib.h>

typedef struct namespace_ {
	char * name;
	char * value;
	char * description;
} namespace;

typedef void (*callbackFn)(namespace ns);

static inline void call_c_func(callbackFn fPtr, namespace ns) {
	(fPtr)(ns);
}
*/
import "C"
import "unsafe"

//export AddDefinition
func AddDefinition(ns string, unit string, isDefault bool, description string) {
}

//export AddFilter
func AddFilter(taskID string, ns string) {
}

//export ProcessMetricWithInt
func ProcessMetricWithInt(taskID string, ns string, v int) {
}

//export Clear
func Clear(taskID string) {
}

//export ListMetrics
//func ListMetrics(taskID string, callback C.callbackFn) {
func ListMetrics(callback C.callbackFn) {
	ns := C.namespace{}
	ns.name = C.CString("--$name")
	ns.value = C.CString("--$value")
	ns.description = C.CString("--$desc")

	C.call_c_func(callback, ns)

	C.free(unsafe.Pointer(ns.name))
	C.free(unsafe.Pointer(ns.value))
	C.free(unsafe.Pointer(ns.description))
}

//export ListMetricsFilters
func ListMetricsFilters(taskID string, callback C.callbackFn) {
}

//export ListMetricsDefinition
func ListMetricsDefinition(callback C.callbackFn) {
}

func main() {}
