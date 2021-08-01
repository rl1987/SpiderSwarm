package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemIsSplayable(t *testing.T) {
	item1 := &Item{
		Fields: map[string]interface{}{
			"field1": "a",
			"field2": "b",
		},
	}

	assert.False(t, item1.IsSplayable())

	item2 := &Item{
		Fields: map[string]interface{}{
			"field1": []interface{}{"a", "b", "c"},
			"field2": []interface{}{"1", "2", "3"},
		},
	}

	assert.True(t, item2.IsSplayable())

	item3 := &Item{
		Fields: map[string]interface{}{
			"field1": []interface{}{"x", "y", "z"},
			"field2": []interface{}{"0", "1"},
		},
	}

	assert.False(t, item3.IsSplayable())

	item4 := &Item{
		Fields: map[string]interface{}{
			"field1": "a",
			"field2": []interface{}{},
		},
	}

	assert.False(t, item4.IsSplayable())
}

func TestItemSplay(t *testing.T) {
	simpleItem := &Item{
		Fields: map[string]interface{}{
			"field1": "a",
			"field2": "1",
		},
	}

	items := simpleItem.Splay()

	assert.Equal(t, []*Item{simpleItem}, items)

	compositeItem := &Item{
		Fields: map[string]interface{}{
			"field1": []interface{}{"a", "b", "c"},
			"field2": []interface{}{"1", "2", "3"},
			"field3": "C",
		},
	}

	expectItems := []*Item{
		&Item{
			Fields: map[string]interface{}{
				"field1": "a",
				"field2": "1",
				"field3": "C",
			},
		},
		&Item{
			Fields: map[string]interface{}{
				"field1": "b",
				"field2": "2",
				"field3": "C",
			},
		},
		&Item{
			Fields: map[string]interface{}{
				"field1": "c",
				"field2": "3",
				"field3": "C",
			},
		},
	}

	items = compositeItem.Splay()

	assert.Equal(t, expectItems[0].Fields, items[0].Fields)
	assert.Equal(t, expectItems[1].Fields, items[1].Fields)
	assert.Equal(t, expectItems[2].Fields, items[2].Fields)
}
