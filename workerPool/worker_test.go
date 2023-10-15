package workerPool

import (
	"context"
	"testing"
)

func TestWorkerInit(t *testing.T) {

	type ExecuteFunc func(val int) (func(), <-chan int)

	testData := []struct {
		name     string
		val      int
		execute  ExecuteFunc
		expected int
	}{
		{
			name: "execution of sum func",
			val:  7,
			execute: func(val int) (func(), <-chan int) {
				ch := make(chan int)
				return func() {

					val += val
					ch <- val
					close(ch)
				}, ch
			},
			expected: 14,
		},
		{
			name: "execution of mulitply func",
			val:  7,
			execute: func(val int) (func(), <-chan int) {
				ch := make(chan int)
				return func() {

					val *= val
					ch <- val
					close(ch)
				}, ch
			},
			expected: 49,
		},
	}

	for _, test := range testData {
		t.Run(
			test.name, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				work := newWorker()

				err := work.Run(ctx)
				if err != nil {
					t.Fatalf("cant start worker: %s", err.Error())
				}
				err = work.Run(ctx)
				if err == nil {
					t.Fatalf("can start worker twice")
				}

				function, ch := test.execute(test.val)

				err = work.SendTask(function)
				if err != nil {
					t.Fatalf("cant send task: %s", err.Error())
				}

				val := <-ch
				if val != test.expected {
					t.Fatalf("val!=test.expected")
				}

				cancel()

				err = work.Close()
				if err != nil {
					t.Fatalf("work.Close() error: %s", err.Error())
				}

				err = work.Close()
				if err == nil {
					t.Fatalf("work.Close() is nil")
				}
			},
		)
	}

}
