package main

import (
	"fmt"
	"time"
)

func main() {
	var bar Bar
	bar.NewOption(0, 100)
	for i := 0; i <= 100; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.play(i)
	}
	bar.finish()
}

type Bar struct {
	percent int
	cur     int
	total   int
	rate    string
	graph   string
}

func (b *Bar) NewOption(start, total int) {
	b.cur = start
	b.total = total
	if b.graph == "" {
		b.graph = "█"
	}
	b.percent = b.getPercent()
	for i := 0; i < b.percent; i = +2 {
		b.rate += b.graph
	}
}

func (b *Bar) getPercent() int {
	return int(float32(b.cur) / float32(b.total) * 100)
}

func (b *Bar) NewOptionWithGraph(start, total int, graph string) {
	b.graph = graph
	b.NewOption(start, total)
}

func (b *Bar) play(cur int) {
	b.cur = cur
	last := b.percent
	b.percent = b.getPercent()
	if b.percent != last && b.percent%2 == 0 {
		b.rate += b.graph
	}

	fmt.Printf("\r[%-50s]%3d%% %8d/%d", b.rate, b.percent, b.cur, b.total)
}

func (b *Bar) finish() {
	fmt.Println()
}
