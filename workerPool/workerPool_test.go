package workerPool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWorkerPoolInit(t *testing.T) {
	t.Run(
		"workerPool init", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			work := NewWorkerPool()

			err := work.Run(ctx)
			if err != nil {
				t.Fatalf("cant run worker pool")
			}

			err = work.Run(ctx)
			if err == nil {
				t.Fatalf("can run worker pool twice")
			}

			work.SendTask(func() {})

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
