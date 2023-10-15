package pipeline_example

import (
	"context"
	"io"

	"github.com/Cladkuu/pipeline/state"
)

type Stage interface {
	Run(ctx context.Context)
	io.Closer
}

type Pipeline struct {
	state  *state.State
	cancel context.CancelFunc
	stages []Stage
}

func NewPipeline(stages ...Stage) *Pipeline {
	p := &Pipeline{state: state.NewState()}

	p.stages = stages

	return p
}

func (p *Pipeline) Run() error {
	if err := p.state.Activate(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	for _, st := range p.stages {
		go st.Run(ctx)
	}

	return nil
}

func (p *Pipeline) Close() error {
	err := p.state.ShutDown()
	if err != nil {
		return err
	}

	p.cancel()

	for _, stage := range p.stages {
		stage.Close()
	}

	err = p.state.Close()

	return err
}
