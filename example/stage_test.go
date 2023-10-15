package example

import (
	"context"
	"fmt"
	pipeline_example "github.com/Cladkuu/pipeline"
	"testing"
	"time"
)

func TestStage(t *testing.T) {
	t.Run(
		"sum and sum", func(t *testing.T) {

			// generator start

			ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			mockGen := NewMockGenerator(ints)

			/*mockGen := make(chan int)
			go func() {
				for _, val := range ints {
					mockGen <- val
				}
				close(mockGen)
			}()*/

			// generator end

			// create stages
			multStage := NewSumStage(mockGen.out)
			multStage2 := NewSumStage(multStage.write)

			// читаем последний канал
			expectedResult := map[int]struct{}{
				0:  struct{}{},
				20: struct{}{},
				8:  struct{}{},
				12: struct{}{},
				24: struct{}{},
				4:  struct{}{},
				16: struct{}{},
				28: struct{}{},
				32: struct{}{},
				36: struct{}{},
				40: struct{}{},
			}
			closeCh := make(chan struct{})
			defer close(closeCh)

			go func() {
				defer func() {
					closeCh <- struct{}{}
				}()

				var ok bool
				for val := range multStage2.write {
					fmt.Println(val)
					_, ok = expectedResult[val]
					if !ok {
						fmt.Printf("get not expectedResult value: %d\n", val)
					}

					delete(expectedResult, val)
				}

				if len(expectedResult) != 0 {
					fmt.Println("len(expectedResult)!=0")
				}

			}()

			pipe := pipeline_example.NewPipeline(mockGen, multStage, multStage2)
			err := pipe.Run(context.Background())
			if err != nil {
				t.Fatalf("pipeline not started: %s", err.Error())
			}

			<-time.After(time.Millisecond * 500)

			err = pipe.Close()
			if err != nil {
				t.Fatalf("pipeline not closed: %s", err.Error())
			}

			<-closeCh
		},
	)
}
