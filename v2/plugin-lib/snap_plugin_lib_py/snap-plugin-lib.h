/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package github.com/librato/snap-plugin-lib-go/v2/plugin-lib */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h> /* for ptrdiff_t below */

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */


#line 3 "main.go"

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
} map_element_t;

typedef struct {
    map_element_t * elements;
    int length;
} map_t;

static inline char * get_map_key(map_t * map, int index) { return map->elements[index].key; }
static inline char * get_map_value(map_t * map, int index) { return map->elements[index].value; }
static inline int get_map_length(map_t * map) { return map->length; }

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
    }
}

typedef struct {
    int sec;
    int nsec;
} time_with_ns_t;

typedef struct {
    map_t * tags_to_add;
    map_t * tags_to_remove;
    time_with_ns_t * timestamp;
    char ** description;
    char ** unit;
} modifiers_t;

static inline modifiers_t * alloc_modifiers() {
    modifiers_t * modifiers = (modifiers_t *) malloc(sizeof(modifiers_t));
    modifiers->tags_to_add = NULL;
    modifiers->tags_to_remove = NULL;
    modifiers->timestamp = NULL;
    modifiers->description = NULL;
    modifiers->unit = NULL;
    return modifiers;
}

static inline void set_modifier_description (modifiers_t * modifiers, char * description) {
    modifiers->description = &description;
}

static inline void set_modifier_unit (modifiers_t * modifiers, char * unit) {
    modifiers->unit = &unit;
}

static inline void set_modifier_tags_to_add (modifiers_t * modifiers, map_t * tags) {
    modifiers->tags_to_add = tags;
}

static inline void set_modifier_tags_to_remove (modifiers_t * modifiers, map_t * tags) {
    modifiers->tags_to_remove = tags;
}

static inline void set_modifier_timestamp (modifiers_t * modifiers, time_with_ns_t * timestamp) {
    modifiers->timestamp = timestamp;
}

static inline char** alloc_str_array(int size) {
    return malloc(sizeof(char*) * size);
}

static inline void set_str_array_element(char **str_array, int index, char *element) {
    str_array[index] = element;
}


#line 1 "cgo-generated-wrapper"


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


extern error_t* ctx_add_metric(char* p0, char* p1, value_t* p2, modifiers_t* p3);

extern error_t* ctx_always_apply(char* p0, char* p1, modifiers_t* p2);

extern void ctx_dismiss_all_modifiers(char* p0);

extern GoInt ctx_should_process(char* p0, char* p1);

extern char** ctx_requested_metrics(char* p0);

extern char* ctx_config(char* p0, char* p1);

extern char** ctx_config_keys(char* p0);

extern char* ctx_raw_config(char* p0);

extern void ctx_add_warning(char* p0, char* p1);

extern GoInt ctx_is_done(char* p0);

extern void ctx_log(char* p0, int p1, char* p2, map_t* p3);

extern void define_metric(char* p0, char* p1, GoInt p2, char* p3);

extern void define_group(char* p0, char* p1);

extern error_t* define_example_config(char* p0);

extern void define_tasks_per_instance_limit(GoInt p0);

extern void define_instances_limit(GoInt p0);

extern void start_collector(callback_t* p0, callback_t* p1, callback_t* p2, define_callback_t* p3, char* p4, char* p5);

#ifdef __cplusplus
}
#endif
