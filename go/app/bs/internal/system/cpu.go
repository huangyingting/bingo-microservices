package system

import (
	bsv1 "bingo/api/bs/v1"
	"runtime"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	TIME_SLICE = 1 * time.Millisecond
)

// CpuLoad represents a CPU load.
type CpuLoad struct {
	cancel chan struct{}
	wg     sync.WaitGroup
	cores  int
	h      *log.Helper
	change []chan int64
}

// NewCpuLoad returns a configured CPU load.
func NewCpuLoad(cores int, h *log.Helper) *CpuLoad {
	return &CpuLoad{
		cores:  cores,
		h:      h,
		change: make([]chan int64, cores),
	}
}

// Start starts up the load goroutines.
func (cl *CpuLoad) Start() {
	cl.cancel = make(chan struct{})
	cl.wg.Add(cl.cores)
	for i := 0; i < cl.cores; i++ {
		ch := make(chan int64, 1)
		cl.change[i] = ch
		go cl.run(i, ch)
	}

}

// Stop signals all goroutines to stop and waits for them to return.
func (cl *CpuLoad) Stop() {
	close(cl.cancel)
	cl.wg.Wait()
}

// Update updates the load percentage of all goroutines.
func (cl *CpuLoad) Update(request *bsv1.CpuLoadRequest) {
	cl.h.Debugf("update cpu load to %d%%\n", request.Percent)
	for _, ch := range cl.change {
		ch <- request.Percent
	}
}

// cpuLoad is a CPU load goroutine.
func (cl *CpuLoad) run(n int, changed <-chan int64) {
	defer cl.wg.Done()

	// Bind the goroutine to an OS thread, so the scheduler won't move it around.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var percent int64 = 0

	// TODO: improve busy loop!
	for {
		begin := time.Now()
		for {
			select {
			case <-cl.cancel:
				return
			case percent = <-changed:
				cl.h.Debugf("update thread %d cpu load to %d%%\n", n, percent)
			default:
				// default branch is required to keep the infinite loop busy.
			}
			if time.Since(begin) > time.Duration(percent)*TIME_SLICE {
				break
			}
		}
		time.Sleep(time.Duration(100-percent) * TIME_SLICE)
	}
}
