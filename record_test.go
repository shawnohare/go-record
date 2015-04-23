package record

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// a sub record of the example above.
var subexample = map[string]interface{}{
	"1": map[string]interface{}{
		"2": 12,
	},
	"3": map[string]interface{}{
		"1": map[string]interface{}{
			"1": map[string]interface{}{
				"1": "value",
			},
		},
	},
}

func TestGet(t *testing.T) {
	var testCases = []struct {
		cmap  map[string]interface{}
		path  string
		value interface{}
		found bool
	}{
		{makeExample(), "1.1", 11, true},
		{makeExample(), "3.1.1.1", "value", true},
		{makeExample(), "4", nil, false},
		{makeExample(), "3.1.2", nil, false},
		{makeExample(), "", nil, false},
	}

	for i, tt := range testCases {
		t.Log("Case:", i)
		r := Init(tt.cmap)
		v, prs := r.Get(tt.path)
		assert.Equal(t, tt.found, prs)
		assert.True(t, reflect.DeepEqual(tt.value, v))
	}

	// Test getting with no path.
	v, prs := get(makeExample(), nil)
	assert.Nil(t, v)
	assert.False(t, prs)
}

func TestSet(t *testing.T) {
	r := Init(makeExample())
	r.Set("1.3", 13)
	s := r.SubRecord([]string{"1"}).Data()
	expected := map[string]interface{}{
		"1": map[string]interface{}{
			"1": 11,
			"2": 12,
			"3": 13,
		},
	}
	assert.Equal(t, expected, s)

	// Test passing in no path to set
	example := makeExample()
	set(example, nil, false)
	assert.Equal(t, example, makeExample())

}

func TestSubRecord(t *testing.T) {
	var testCases = []struct {
		example    map[string]interface{}
		paths      []string
		subexample map[string]interface{}
	}{
		{
			makeExample(),
			[]string{"1.2", "3"},
			subexample,
		},
		// Test whether SubRecord ignores invalid paths.
		{
			makeExample(),
			[]string{"1.2", "3", "4.1.2.3"},
			subexample,
		},
	}

	for i, tt := range testCases {
		t.Log("Case:", i)
		r := Init(tt.example)
		sub2 := r.SubRecord(tt.paths).Data()
		assert.True(t, reflect.DeepEqual(tt.subexample, sub2))
	}

	// Test passing all paths to some leafs vs path to parent.
	r := Init(makeExample())
	sParent := r.SubRecord([]string{"1"})
	sAll := r.SubRecord([]string{"1.1", "1.2"})
	assert.Equal(t, sParent, sAll)
}
