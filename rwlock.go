// Reference: https://eli.thegreenplace.net/2019/implementing-reader-writer-locks
package goexamples

type RWLocker interface {
	RLock()
	RUnlock()
	WLock()
	WUnlock()
}
