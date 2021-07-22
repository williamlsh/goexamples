package goexamples

import (
	"sync"
	"sync/atomic"
)

// WritePreferFastRWLock is a faster implementation of a RW lock that prefers
// writers. It's similar to the Go standard library implementation of RWMutex,
// but uses channels instead of runtime_Sem* primitives.
type WriterPreferFastRWLock struct {
	// w guarantees mutual exclusion between writers.
	w sync.Mutex

	// These channels serve as semaphores for a writer to wait on readers and for
	// readers to wait for a writer, without spinning.
	writerWait chan struct{}
	readerWait chan struct{}

	// numPending marks the number of readers or writers that are still using the
	// lock. Readers increment it by one, but a writer subtracts maxReaders;
	// therefore, a negative value means a writer is currently using the lock.
	numPending int32

	// readersDeparting is used when a writer waits for all existing readers to
	// flush out before doing its thing. When a writer comes in to take a look
	// and guarantees it's the only writer, it posts -maxReaders onto numPending
	// to notify new readers not to enter. It then also posts the number of
	// pending readers onto readersDeparting. Each reader will decrement this
	// field before relinquishing the lock, and the writer waits for all of them
	// before proceeding.
	readersDeparting int32
}

const maxReaders int32 = 1 << 30

func NewWriterPreferFastRWLock() *WriterPreferFastRWLock {
	var l WriterPreferFastRWLock
	l.writerWait = make(chan struct{})
	l.readerWait = make(chan struct{})
	return &l
}

func (l *WriterPreferFastRWLock) RLock() {
	if atomic.AddInt32(&l.numPending, 1) < 0 {
		// A writer is pending, wait for it.
		<-l.readerWait
	}
}

func (l *WriterPreferFastRWLock) RUnlock() {
	if r := atomic.AddInt32(&l.numPending, -1); r < 0 {
		// If numPending is now negative, it can be either because of an
		// error ...
		if r+1 == 0 || r+1 == -maxReaders {
			panic("RUnlock of unlocked RWLock")
		}
		// ... or because a writer is pending
		if atomic.AddInt32(&l.readersDeparting, -1) == 0 {
			l.writerWait <- struct{}{}
		}
	}
}

func (l *WriterPreferFastRWLock) WLock() {
	l.w.Lock()
	// Announce to readers there is a pending writer by decrementing maxReaders
	// from numPending. r will hold the number of pending readers.
	r := atomic.AddInt32(&l.numPending, -maxReaders) + maxReaders
	if r != 0 && atomic.AddInt32(&l.readersDeparting, r) != 0 {
		<-l.readerWait
	}
}

func (l *WriterPreferFastRWLock) WUnlock() {
	// Announce to readers there is no longer an active writer.
	r := atomic.AddInt32(&l.numPending, maxReaders)
	if r >= maxReaders {
		panic("WUnlock of unlocked RWLock")
	}
	// Unblock blocked readers, if any.
	for i := 0; i < int(r); i++ {
		l.readerWait <- struct{}{}
	}
	// Allow other writers to proceed.
	l.w.Unlock()
}
