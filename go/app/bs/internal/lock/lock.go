package lock

type IDistributedLock interface {
	Lock(pfx string) error
	Unlock() error
	Close()
}
