package pipeline_example

import (
	"context"
	"fmt"
	"io"

	"github.com/Cladkuu/pipeline/state"
)

type Generator interface {
	Generate(ctx context.Context)
	io.Closer
}

type Stage interface {
	Run(ctx context.Context)
	io.Closer
}

type Pipeline struct {
	generator Generator
	state     *state.State
	cancel    context.CancelFunc
	stages    []Stage
}

func NewPipeline(generator Generator, stages ...Stage) *Pipeline {
	p := &Pipeline{
		state:     state.NewState(),
		generator: generator,
	}

	p.stages = stages

	return p
}

func (p *Pipeline) Run(ctx context.Context) error {
	if err := p.state.Activate(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	p.generator.Generate(ctx)

	for _, st := range p.stages {
		st.Run(ctx)
	}

	return nil
}

func (p *Pipeline) Close() error {
	err := p.state.ShutDown()
	if err != nil {
		return err
	}

	p.cancel()

	p.generator.Close()

	for _, stage := range p.stages {
		stage.Close()
	}

	err = p.state.Close()
	fmt.Println("pipe closed")

	return err
}
