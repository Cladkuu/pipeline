package pipeline_example

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {

	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := &st{
		ints:  ints,
		write: make(chan int),
	}

	pipe := NewPipeline(nil, s)
	err := pipe.Run()
	if err != nil {
		t.Fatalf("pipe.Run err: %s", err.Error())
	}

	done := make(chan int)
	go func() {
		time.Sleep(time.Millisecond * 250)
		close(done)
	}()

	defer func() {
		err = pipe.Close()
		if err != nil {
			t.Fatalf("pipe.Close err: %s", err.Error())
		}
	}()

	for {
		select {
		case val := <-s.write:
			fmt.Println(val)
		case <-done:
			return
		}
	}

}

type st struct {
	ints  []int
	write chan int
}

func (s *st) Run(ctx context.Context) {

	go func() {
		for _, val := range s.ints {
			time.Sleep(time.Millisecond * 100)
			s.write <- val
		}
		fmt.Println("closing")

		close(s.write)
	}()

}

func (s *st) Close() error {
	return nil
}
