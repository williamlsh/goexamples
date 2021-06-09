package x

import (
	"fmt"
	"reflect"
	"testing"
)

func TestProcessQuery(t *testing.T) {
	var cache Cache
	t.Run("empty query", func(t *testing.T) {
		var content []string
		if result := processQuery(content, &cache); !reflect.DeepEqual(result, []bool{}) {
			t.Fatalf("unexpected result: %+v", result)
		}
		fmt.Printf("cache in test 1: %+v\n", cache.list)
	})
	t.Run("Two times query one hit", func(t *testing.T) {
		content1 := []string{"a", "b", "c"}
		if result := processQuery(content1, &cache); !reflect.DeepEqual(result, []bool{false, false, false}) {
			t.Fatalf("unexpected result: %+v", result)
		}

		content2 := []string{"a", "c", "d"}
		if result := processQuery(content2, &cache); !reflect.DeepEqual(result, []bool{true, false, false}) {
			t.Fatalf("unexpected result: %+v", result)
		}

		if !reflect.DeepEqual(cache.list, []Set{
			{"a": struct{}{}},
			{
				"b": struct{}{},
				"c": struct{}{},
			},
			{
				"c": struct{}{},
				"d": struct{}{},
			},
		}) {
			t.Fatalf("Unexpected cache result: %+v\n", cache.list)
		}
	})
	t.Run("shorter query content all hits", func(t *testing.T) {
		content := []string{"a", "b"}
		if result := processQuery(content, &cache); !reflect.DeepEqual(result, []bool{true, true}) {
			t.Fatalf("unexpected result: %+v", result)
		}
	})
	t.Run("longger query content no hit", func(t *testing.T) {
		content := []string{"e", "f", "g", "h", "i"}
		if result := processQuery(content, &cache); !reflect.DeepEqual(result, []bool{false, false, false, false, false}) {
			t.Fatalf("unexpected result: %+v", result)
		}
		fmt.Printf("cache: %+v\n", cache.list)
	})
}
