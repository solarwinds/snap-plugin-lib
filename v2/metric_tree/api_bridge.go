package main

import "C"

/*
#include <stdlib.h>

typedef struct {
	char * name;
	char * value;
	char * description;
} namespaceElement;

typedef struct {
	long long length;
	namespaceElement *elements;
} namespace;

static inline namespace * allocNamespace(long long len) {
	namespace * ns = malloc(sizeof(namespace));

	namespaceElement * nsEls = malloc(sizeof(namespaceElement) * len);
	ns->length = len;
	ns->elements = nsEls;

	return ns;
}

static inline void freeNamespace(namespace * ns) {
	for (int i = 0; i < ns->length; ++i) {
		free(ns->elements[i].name);
		free(ns->elements[i].value);
		free(ns->elements[i].description);
	}
	free(ns->elements);
	free(ns);
}

static inline void setNamespaceElement(namespace * ns, int el, char * name, char * value, char * description) {
	ns->elements[el].name = name;
	ns->elements[el].value = value;
	ns->elements[el].description = description;
}

enum metricValueType {
	METRIC_VALUE_INT,
	METRIC_VALUE_DOUBLE
};

typedef union {
	int intValue;
	double doubleValue;
} metricValueData;

typedef struct {
	int type;
	metricValueData data;
} metricValue;

static inline metricValue* withIntData(int data) {
	metricValue *mv = malloc(sizeof(metricValue));
	mv->type = METRIC_VALUE_INT;
	mv->data.intValue = data;
	return mv;
}

typedef void (*callbackFn)(namespace * ns);

static inline void call_c_callback(callbackFn fPtr, namespace * ns) {
	(fPtr)(ns);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

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
	nss := C.allocNamespace(4)

	nssPtr := (*C.namespace)(unsafe.Pointer(nss))
	C.setNamespaceElement(nssPtr, 0, C.CString("name111"), C.CString("value111"), C.CString("desc111"))
	C.setNamespaceElement(nssPtr, 1, C.CString("name222"), C.CString("value222"), C.CString("desc222"))
	C.setNamespaceElement(nssPtr, 2, C.CString("name333"), C.CString("value333"), C.CString("desc333"))
	C.setNamespaceElement(nssPtr, 3, C.CString("name444"), C.CString("value444"), C.CString("desc444"))

	data := C.withIntData(10)
	dataPtr := (*C.metricValue)(unsafe.Pointer(data))
	fmt.Printf("dataPtr=%#v\n", dataPtr)

	C.call_c_callback(callback, nssPtr)
	C.freeNamespace(nssPtr)
}

//export ListMetricsFilters
func ListMetricsFilters(taskID string, callback C.callbackFn) {
}

//export ListMetricsDefinition
func ListMetricsDefinition(callback C.callbackFn) {
}

func main() {}
