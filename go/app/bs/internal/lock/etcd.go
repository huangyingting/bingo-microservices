package lock

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
)

var ErrInvalidEtcdAddresses error = errors.New("invalid etcd addresses")

type EtcdDistributedLock struct {
	cli     *clientv3.Client
	session *concurrency.Session
	lock    *concurrency.Mutex
	h       *log.Helper
}

func NewEtcdDistributedLock(addr []string, h *log.Helper) (*EtcdDistributedLock, error) {
	if len(addr) == 0 {
		return nil, ErrInvalidEtcdAddresses
	}
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   addr,
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
			DialTimeout: 5 * time.Second})

	if err != nil {
		h.Errorf("connect to etcd failed: %v", err)
		return nil, err
	}

	session, err := concurrency.NewSession(cli)
	if err != nil {
		h.Errorf("open etcd session failed: %v", err)
		cli.Close()
		return nil, err
	}

	return &EtcdDistributedLock{cli: cli, session: session, h: h}, nil
}

func (l *EtcdDistributedLock) Lock(pfx string) error {
	l.lock = concurrency.NewMutex(l.session, pfx)
	err := l.lock.TryLock(context.Background())
	if err != nil {
		l.h.Errorf("lock error: %v", err)
		l.lock = nil
	}
	return err
}

func (l *EtcdDistributedLock) Unlock() error {
	if l.lock != nil {
		err := l.lock.Unlock(context.Background())
		if err != nil {
			l.h.Errorf("unlock error: %v", err)
		}
		return err
	}
	return nil
}

func (l *EtcdDistributedLock) Close() {
	l.Unlock()

	if l.session != nil {
		l.session.Close()
	}

	if l.cli != nil {
		l.cli.Close()
	}
}
