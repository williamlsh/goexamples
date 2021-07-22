package goexamples

import "sync"

type ReaderCountRWLock struct {
	m           sync.Mutex
	readerCount int
}

func (l *ReaderCountRWLock) RLock() {
	l.m.Lock()
	l.readerCount++
	l.m.Unlock()
}

func (l *ReaderCountRWLock) RUnlock() {
	l.m.Lock()
	l.readerCount--
	l.m.Unlock()
}

func (l *ReaderCountRWLock) WLock() {
	for {
		l.m.Lock()
		if l.readerCount > 0 {
			l.m.Unlock()
		} else {
			break
		}
	}
}

func (l *ReaderCountRWLock) WUnlock() {
	l.m.Unlock()
}
