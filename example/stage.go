package example

import (
	"context"
	"fmt"
	"github.com/Cladkuu/pipeline/workerPool"
)

type sumStage struct {
	read       chan int
	write      chan int
	closeCH    chan struct{}
	cancel     context.CancelFunc
	workerPool *workerPool.WorkerPool
}

func NewSumStage(
	read chan int,
) *sumStage {
	return &sumStage{
		closeCH:    make(chan struct{}),
		read:       read,
		write:      make(chan int),
		workerPool: workerPool.NewWorkerPool(),
	}
}

func (s *sumStage) Run(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)

	s.workerPool.Run(ctx)

	go func() {
		defer func() {
			s.closeCH <- struct{}{}
		}()

		var err error

		for {
			select {
			case val, ok := <-s.read:
				if !ok {
					return
				}
				err = s.workerPool.SendTask(s.createExecution(val))
				if err != nil {
					fmt.Printf("err occured: %s\n", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

}

func (s *sumStage) createExecution(val int) workerPool.Request {
	return func() {
		val += val

		s.write <- val
	}
}

func (s *sumStage) Close() error {
	s.cancel()
	// дожидаемся завершения приема обращений
	<-s.closeCH
	close(s.closeCH)

	fmt.Println("STAGE STARTING CLOSING")

	// завершаем обработку всех принятых обращений
	s.workerPool.Close()

	// закрываем канал для отправки обработанных сообщений
	close(s.write)

	return nil
}
