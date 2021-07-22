package goexamples

import (
	"context"

	"golang.org/x/sync/semaphore"
)

const maxWeight int64 = 1 << 30

type SemaRWLock struct {
	s *semaphore.Weighted
}

func NewSemaRWLock() *SemaRWLock {
	return &SemaRWLock{
		semaphore.NewWeighted(maxWeight),
	}
}

func (l *SemaRWLock) RLock() {
	l.s.Acquire(context.TODO(), 1)
}

func (l *SemaRWLock) RUnlock() {
	l.s.Release(1)
}

func (l *SemaRWLock) WLock() {
	l.s.Acquire(context.TODO(), maxWeight)
}

func (l *SemaRWLock) WUnlock() {
	l.s.Release(maxWeight)
}
