package workerPool

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Cladkuu/pipeline/state"
)

type Request func()

type Worker interface {
	io.Closer
	SendTask(req Request) error
	Run(ctx context.Context) error
}

type workers struct {
	workerCount int
	workers     []Worker
	workerIndex int
	mutex       *sync.Mutex
}

func (w *workers) Close() error {
	for _, work := range w.workers {
		work.Close()
	}

	return nil
}

func (w *workers) Run(ctx context.Context) error {
	var err error
	for _, work := range w.workers {
		err = work.Run(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *workers) GetWorker() Worker {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.workerIndex++
	if w.workerIndex == len(w.workers) {
		w.workerIndex = 0
	}

	return w.workers[w.workerIndex]
}

func newWorkers() *workers {
	w := &workers{
		mutex:       &sync.Mutex{},
		workerCount: 5,
	}

	w.workers = make([]Worker, 0, w.workerCount)
	for i := 0; i < w.workerCount; i++ {
		w.workers = append(
			w.workers, newWorker(),
		)
	}

	return w
}

type WorkerPool struct {
	requests <-chan Request
	workers  *workers
	cancel   context.CancelFunc
	state    *state.State
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		workers:  newWorkers(),
		state:    state.NewState(),
		requests: make(chan Request),
	}
}

func (w *WorkerPool) Run(ctx context.Context) error {
	err := w.state.Activate()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	err = w.workers.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (w *WorkerPool) SendTask(req Request) error {
	if err := w.state.IsActive(); err != nil {
		return err
	}

	return w.workers.GetWorker().SendTask(req)
}

func (w *WorkerPool) Close() error {
	err := w.state.ShutDown()
	if err != nil {
		return err
	}

	fmt.Println("WORKER POOL STARTING CLOSIng")
	w.cancel()

	w.workers.Close()

	err = w.state.Close()

	fmt.Println("worker_pool closed")

	return err
}
