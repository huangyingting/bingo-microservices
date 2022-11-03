package system

import (
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	bsv1 "bingo/api/bs/v1"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	MB_IN_BYTES = 1 << 20
)

// MemLoad represents a memory load.
type MemLoad struct {
	cancel   chan struct{}
	wg       sync.WaitGroup
	alloc    [][]byte
	change   chan int
	pageSize int
	h        *log.Helper
}

// NewLoad returns a configured memory load.
func NewMemLoad(h *log.Helper) *MemLoad {
	return &MemLoad{
		change:   make(chan int, 1),
		pageSize: os.Getpagesize(),
		h:        h,
	}
}

// Start starts up the load goroutine.
func (ml *MemLoad) Start() {
	ml.cancel = make(chan struct{})

	ml.wg.Add(1)
	go ml.run()
}

// Stop signals the load goroutine to stop and waits for it to return.
func (ml *MemLoad) Stop() {
	close(ml.cancel)
	ml.wg.Wait()
}

// Update updates the allocated memory size.
func (ml *MemLoad) Update(request *bsv1.MemLoadRequest) {
	ml.h.Debugf("update memory load to %dMB", request.Size)
	ml.change <- int(request.Size)
}

func (ml *MemLoad) run() {
	defer ml.wg.Done()
	defer func() {
		ml.alloc = nil
		runtime.GC()
	}()

	for {
		// do not use default branch in select as we don't want a busy loop.
		select {
		case <-ml.cancel:
			return
		case size := <-ml.change:
			ml.alloc = nil
			runtime.GC()
			for page := 0; page < size*MB_IN_BYTES/ml.pageSize; page++ {
				// allocate memory in page-sized chunks.
				chunk := make([]byte, ml.pageSize)
				ml.alloc = append(ml.alloc, chunk)
			}
			ml.h.Debugf(
				"memory - page size: %d bytes, pages: %d, size: %d MB\n",
				ml.pageSize,
				len(ml.alloc),
				len(ml.alloc)*ml.pageSize/MB_IN_BYTES,
			)
		case <-time.After(time.Second):
			// make sure we use the allocated memory, so it won't get swapped.
			if ml.alloc != nil {
				for page := 0; page < len(ml.alloc); page++ {
					ml.alloc[page][rand.Intn(ml.pageSize)]++
				}
			}
		}
	}
}
