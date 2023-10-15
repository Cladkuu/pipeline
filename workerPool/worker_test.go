package workerPool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWorkerInit(t *testing.T) {

	t.Run(
		"worker init", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			work := newWorker()

			err := work.Run(ctx)
			if err != nil {
				t.Fatalf("cant start worker")
			}
			err = work.Run(ctx)
			if err == nil {
				t.Fatalf("can start worker twice")
			}

			work.SendTask(func() {})

			cancel()

			time.Sleep(time.Second * 2)

			err = work.Close()
			if err != nil {
				t.Fatalf("work.Close() error: %s", err.Error())
			}
			fmt.Println(err)

			err = work.Close()
			if err == nil {
				t.Fatalf("work.Close() is nil")
			}
			fmt.Println(err)

		},
	)

}
