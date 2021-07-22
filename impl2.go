package goexamples

import "sync"

type ReaderCountCondRWLock struct {
	readerCount int
	c           *sync.Cond
}

func NewReaderCountCondRWLock() *ReaderCountCondRWLock {
	return &ReaderCountCondRWLock{
		0,
		sync.NewCond(new(sync.Mutex)),
	}
}

func (l *ReaderCountCondRWLock) RLock() {
	l.c.L.Lock()
	l.readerCount++
	l.c.L.Unlock()
}

func (l *ReaderCountCondRWLock) RUnlock() {
	l.c.L.Lock()
	l.readerCount--
	if l.readerCount == 0 {
		l.c.Signal()
	}
	l.c.L.Unlock()
}

func (l *ReaderCountCondRWLock) WLock() {
	l.c.L.Lock()
	for l.readerCount > 0 {
		l.c.Wait()
	}
}

func (l *ReaderCountCondRWLock) WUnlock() {
	l.c.Signal()
	l.c.L.Unlock()
}
