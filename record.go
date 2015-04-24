// Package record provides a thin wrapper over a native composite map
// data structure that simplifies access to the final (leaf) values.
// Some database SDKs return such composite map structures
// when fetching a record.
//
// - We assume all keys are strings in order to simplify access.
// - Record values are accessed/set by providing a string representation
// of the path separated by ".". For example, passing "key1.key2" to a Record
// fetchs the value with key "key2" in the subrecord S of the record R, where the key
// of S is "key1".
//
// Mathematically, we can view a composite map as a function
// f whose domain is the subset X of finite length tuples of
// the countable product of some key set A.  For each index i,
// fi is a function from A to the set of possible values.  Then
// f(k1, k2, k3, ..., kn) := f1 (f2 ... fn(kn)).  This makes the
// dot notation for paths quite natural.
package record

import (
	"strings"
)

type Record struct {
	data map[string]interface{}
}

// New creates a pointer to a Record which is ready to have values set.
func New() *Record {
	return &Record{make(map[string]interface{})}
}

// Init creates a pointer to a Record whose underlying data is the input.
func Init(m map[string]interface{}) *Record {
	r := New()
	r.data = m
	return r
}

// AsMap returns the underlying map-of-maps data structure a Record
func (r *Record) AsMap() map[string]interface{} {
	return r.data
}

// Get the element associated to the path.
func (r *Record) Get(path string) (interface{}, bool) {
	p := strings.Split(path, ".")
	return get(r.data, p)
}

// Set inserts the input into the Record.
func (r *Record) Set(path string, x interface{}) {
	p := strings.Split(path, ".")
	set(r.data, p, x)
}

// Filter returns a new Record that only includes the specified paths.
func (r *Record) Filter(paths []string) *Record {
	data := filter(r.data, paths)
	return Init(data)
}

// FilterMap returns a new composite map filtered to include only values
// (and nested maps) specified by the paths array.
func FilterMap(m map[string]interface{}, paths []string) map[string]interface{} {
	return filter(m, paths)
}

// SubRecord produces a new nested-map structure from the input
func filter(m map[string]interface{}, paths []string) map[string]interface{} {
	// Create empty composite map.
	s := make(map[string]interface{})
	for _, pathStr := range paths {
		// Don't insert any value into s if the path in m doesn't exist.
		p := strings.Split(pathStr, ".")
		if v, prs := get(m, p); prs {
			set(s, p, v)
		}
	}
	return s
}

func get(m map[string]interface{}, path []string) (interface{}, bool) {
	l := len(path)
	switch l {
	case 0:
		return nil, false
	case 1:
		x, prs := m[path[0]]
		return x, prs
	default:
		m2 := m[path[0]]
		switch m2.(type) {
		case map[string]interface{}:
			return get(m2.(map[string]interface{}), path[1:])
		default:
			// Invalid key.
			return nil, false
		}
	}
}

func set(m map[string]interface{}, path []string, x interface{}) {
	l := len(path)
	switch l {
	case 0:
		return
	case 1:
		m[path[0]] = x
	default:
		if _, prs := m[path[0]]; !prs {
			m[path[0]] = make(map[string]interface{})
		}
		set(m[path[0]].(map[string]interface{}), path[1:], x)
	}
}

// NOTE the following functions are not being used for anything.

// func getOld(x interface{}, path []string) (interface{}, bool) {
// 	l := len(path)
// 	switch l {
// 	case 0:
// 		return x, true
// 	default:
// 		switch x.(type) {
// 		case map[string]interface{}:
// 			key := path[0]
// 			newMap := x.(map[string]interface{})[key]
// 			return getOld(newMap, path[1:])
// 		default:
// 			// Invalid path.
// 			return nil, false
// 		}
// 	}
// }
