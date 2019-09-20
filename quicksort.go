package main

import (
	"fmt"
	"strings"
)

const dot = "."

func main() {
	ints := []int{-6, 4, -3, 5, -2, -1, 0, 1, -9}
	v := sort(ints)
	fmt.Println(v)

	dict := map[string]int{"A": 1, "B.A": 2, "B.B": 3, "CC.D.E": 4, "CC.D.F": 5}
	vv := transform(dict)
	fmt.Println(vv)
}

func sort(v []int) []int {
	j := 0
	for i, e := range v {
		if e >= 0 {
			continue
		}
		if i != j {
			v[i] = v[j]
			v[j] = e
		}
		j++
	}
	return v
}

func transform(m map[string]int) map[string]interface{} {
	container := make(map[string]interface{})
	for k, v := range m {
		extract(k, v, container)
	}
	return container
}

func extract(k string, v int, container map[string]interface{}) {
	if !strings.Contains(k, dot) {
		container[k] = v
		return
	}
	ks := strings.SplitN(k, dot, 2)
	if _, ok := container[ks[0]]; !ok {
		child := make(map[string]interface{})
		container[ks[0]] = child
	}
	extract(ks[1], v, container[ks[0]].(map[string]interface{}))
}
