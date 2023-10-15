package workerPool

import (
	"context"
	"fmt"

	"github.com/Cladkuu/pipeline/state"
)

type worker struct {
	requests chan Request
	closeCh  chan struct{}
	state    *state.State
}

func newWorker() *worker {
	w := &worker{
		closeCh:  make(chan struct{}),
		requests: make(chan Request),
		state:    state.NewState(),
	}

	return w
}

func (w *worker) Run(ctx context.Context) error {
	if err := w.state.Activate(); err != nil {
		return err
	}

	go w.work(ctx)

	return nil
}

func (w *worker) work(ctx context.Context) {
	defer func() {
		w.closeCh <- struct{}{}
	}()

	for {
		select {
		case req, ok := <-w.requests:
			if !ok {
				return
			}

			req()

		case <-ctx.Done():
			return
		}
	}
}

func (w *worker) SendTask(req Request) error {
	if err := w.state.IsActive(); err != nil {
		return err
	}

	w.requests <- req
	return nil
}

func (w *worker) Close() error {
	err := w.state.ShutDown()
	if err != nil {
		return err
	}

	<-w.closeCh
	fmt.Println("worker closed")
	close(w.closeCh)
	close(w.requests)

	err = w.state.Close()

	return err
}
