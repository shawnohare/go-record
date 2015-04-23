# go-record
A simple Record structure that simplifies access to values in composite maps.

## Usage Examples

The code for these examples is in the [examples](examples/) dir.  We first import
the package and define a function that records a sample raw record
such as might be returned by a database SDK.

```go
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
```

### New 
To create a new empty record, use New.
```go
	var r *record.Record
	r = record.New()
```

### Init
We can also initialize a Record from an existing composite map.
 ```go
	cmap := makeExample()
	r = record.Init(cmap)
```

### Set and Get
We can insert a value into the record with Set.
```go
	r.Set("1.3", 13)
```
This actually maps "3" to 13 in the map keyed to "1".
The value 13 now resides as a value in a subsubmap of the record.
We can extract it with Get, although we must assert types.
```go
	x, _ := r.Get("1.3")
	fmt.Println("Fetched the value 13:", x.(int) == 13) // true
```
This actually maps "3" to 13 in the map keyed to "1".

Get's second return value indicates whether the path exists.
```go
	_, prs := r.Get("1.4")
	if !prs {
		fmt.Println("The path \"1.4\" does not exist.")
	}

```
If we wish to obtain a nested map instead of a leaf, we
just pass the appropriate path.
```go
	subMap, _ := r.Get("1")
	fmt.Println("A nested map:", subMap)
```

### AsMap
 AsMap allows access to the underlying composite map.
```go
	var d map[string]interface{}
	d = r.AsMap()
	fmt.Println("The underlying composite map:", d)
```

### Filter
Similarly, we can retrieve a filtered version of the record
by passing the desired paths. If we pass a path to a non-leaf node,
we obtain all values below the path node as well. Filter
silently ignores  non-existent paths.
```go
	paths := []string{"1.1", "3.1", "badPath"}
	var filtered *record.Record
	filtered = r.Filter(paths)
	fmt.Println("Filtered record:", filtered)
	_, prs = filtered.Get("3.1.1.1") // Access the leaf as expected.
	fmt.Println("Value exists:", prs)
}
```
