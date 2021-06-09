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
		result := processQuery(content, &cache)
		if !reflect.DeepEqual(result, []bool{}) {
			t.Fatalf("unexpected result: %+v", result)
		}
		fmt.Printf("cache in test 1: %+v\n", cache.Set)
	})
	t.Run("Two times query one hit", func(t *testing.T) {
		content1 := []string{"a", "b", "c"}
		result1 := processQuery(content1, &cache)
		if !reflect.DeepEqual(result1, []bool{false, false, false}) {
			t.Fatalf("unexpected result: %+v", result1)
		}

		content2 := []string{"a", "c", "d"}
		result2 := processQuery(content2, &cache)
		if !reflect.DeepEqual(result2, []bool{true, true, false}) {
			t.Fatalf("unexpected result: %+v", result2)
		}

		if !reflect.DeepEqual(cache.Set, Set{
			"a": struct{}{},
			"b": struct{}{},
			"c": struct{}{},
			"d": struct{}{},
		}) {
			t.Fatalf("Unexpected cache result: %+v\n", cache.Set)
		}
	})
	t.Run("Two times query no hit", func(t *testing.T) {
		content1 := []string{"x", "y", "z"}
		result1 := processQuery(content1, &cache)
		if !reflect.DeepEqual(result1, []bool{false, false, false}) {
			t.Fatalf("unexpected result: %+v", result1)
		}

		content2 := []string{"p", "q", "v"}
		result2 := processQuery(content2, &cache)
		if !reflect.DeepEqual(result2, []bool{false, false, false}) {
			t.Fatalf("unexpected result: %+v", result2)
		}
	})
}
