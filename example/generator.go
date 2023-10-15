package example

import (
	"context"
	"fmt"
	"time"
)

type mockGenerator struct {
	// some dependencies: repositories, services and etc
	ints    []int
	out     chan int
	closeCH chan struct{}
}

func NewMockGenerator(ints []int) *mockGenerator {
	return &mockGenerator{
		ints:    ints,
		out:     make(chan int),
		closeCH: make(chan struct{}),
	}
}

func (m *mockGenerator) Generate(ctx context.Context) {
	go func() {

		defer func() {
			close(m.out)
			m.closeCH <- struct{}{}
		}()

		i := 0
		for i < len(m.ints) {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(time.Millisecond * 100):
				m.out <- m.ints[i]
				fmt.Println("generated i")
				i++
			}
		}
	}()
}

func (m *mockGenerator) Close() error {

	<-m.closeCH
	close(m.closeCH)
	fmt.Println("GENERATOR CLOSED!!!!!!!")

	return nil
}
