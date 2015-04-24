package main

import (
	"fmt"

	"github.com/shawnohare/go-record"
)

//
func makeExample() map[string]interface{} {
	var example = map[string]interface{}{
		"1": map[string]interface{}{
			"1": 11,
			"2": 12,
		},
		"2": map[string]interface{}{
			"1": 21,
			"2": 22,
		},
		"3": map[string]interface{}{
			"1": map[string]interface{}{
				"1": map[string]interface{}{
					"1": "value",
				},
			},
		},
	}

	return example
}

func main() {
	var r *record.Record
	// To create a new empty record, use New.
	r = record.New()

	// We can also initialize a Record from an existing composite map.
	// First generate a fresh composite map.
	cmap := makeExample()

	// To Create a new Record wrapping this composite map, we use Init.
	r = record.Init(cmap)

	// We can insert a value into the record with Set.
	r.Set("1.3", 13)

	// The value 13 now resides as a value in a subsubmap of the record.
	// We can extract it with Get, although we must assert types.
	x, _ := r.Get("1.3")
	fmt.Println("Fetched the value 13:", x.(int) == 13) // true

	// Get's second return value indicates whether the path exists.
	_, prs := r.Get("1.4")
	if !prs {
		fmt.Println("The path \"1.4\" does not exist.")
	}

	// If we wish to obtain a nested map instead of a leaf, we
	// just pass the appropriate path.
	subMap, _ := r.Get("1")
	fmt.Println("A nested map:", subMap)

	// AsMap allows access to the underlying composite map.
	var d map[string]interface{}
	d = r.AsMap()
	fmt.Println("The underlying composite map:", d)

	// Similarly, we can retrieve a filtered version of the record
	// by passing the desired paths. If we pass a path to a non-leaf node,
	// we obtain all values below the node as well. Filter
	// silently ignores  non-existent paths.
	paths := []string{"1.1", "3.1", "badPath"}
	var filtered *record.Record
	filtered = r.Filter(paths)
	fmt.Println("Filtered record:", filtered)
	_, prs = filtered.Get("3.1.1.1") // Access the leaf as expected.
	fmt.Println("Value exists:", prs)

	var filteredMap map[string]interface{}
	filteredMap = record.FilterMap(makeExample(), paths)
	fmt.Println("Filtered map:", filteredMap)
	// Equivalent to:
	filteredMap = record.Init(makeExample()).Filter(paths).AsMap()
	fmt.Println("Filtered map:", filteredMap)
}
