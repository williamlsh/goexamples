package goexamples

import "sync"

type WriterPreferRWLock struct {
	readerCount int
	hasWriter   bool
	c           *sync.Cond
}

func NewWriterPreferRWLock() *WriterPreferRWLock {
	return &WriterPreferRWLock{
		0,
		false,
		sync.NewCond(new(sync.Mutex)),
	}
}

func (l *WriterPreferRWLock) RLock() {
	l.c.L.Lock()
	for l.hasWriter {
		l.c.Wait()
	}
	l.readerCount++
	l.c.L.Unlock()
}

func (l *WriterPreferRWLock) RUnlock() {
	l.c.L.Lock()
	l.readerCount--
	if l.readerCount == 0 {
		l.c.Broadcast()
	}
	l.c.L.Unlock()
}

func (l *WriterPreferRWLock) WLock() {
	l.c.L.Lock()
	for l.hasWriter {
		l.c.Wait()
	}
	l.hasWriter = true
	for l.readerCount > 0 {
		l.c.Wait()
	}
	l.c.L.Unlock()
}

func (l *WriterPreferRWLock) WUnlock() {
	l.c.L.Lock()
	l.hasWriter = false
	l.c.Broadcast()
	l.c.L.Unlock()
}
