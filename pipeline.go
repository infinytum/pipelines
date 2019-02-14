package pipelines

import (
	"errors"
	"reflect"

	"github.com/infinytum/reactive"
)

// Pipeline is a special kind of Subject
type Pipeline struct {
	subject   reactive.Subjectable
	pipestack reactive.Observable
}

// Add a new pipe to this pipeline
func (p *Pipeline) Add(pipe interface{}) error {
	if pipe == nil || reflect.TypeOf(pipe).Kind() != reflect.Func {
		return errors.New("pipe is not a function")
	}

	pipedSubject := p.pipestack.Pipe(
		func(parent reactive.Observable, next reactive.Subjectable) {
			parent.Subscribe(func(args ...interface{}) {
				returnVals := notifySubscriber(pipe, args)
				next.Next(returnVals...)
			})
		},
	)
	p.pipestack = pipedSubject.(reactive.Subjectable)
	return nil
}

func (p *Pipeline) Run(values ...interface{}) reactive.Observable {
	sub := reactive.NewReplaySubject()
	p.pipestack.Pipe(reactive.Take(1)).Subscribe(func(args ...interface{}) {
		sub.Next(args...)
	})
	p.subject.Next(values...)
	return sub
}

// NewPipeline returns a default pipeline without buffer
func NewPipeline() *Pipeline {
	sub := reactive.NewSubject()
	return &Pipeline{
		subject:   sub,
		pipestack: sub,
	}
}

// NewReplayPipeline returns a pipeline with a buffer of 1
func NewReplayPipeline() *Pipeline {
	sub := reactive.NewReplaySubject()
	return &Pipeline{
		subject:   sub,
		pipestack: sub,
	}
}

// NewBufferedPipeline returns a pipeline with a custom buffer
func NewBufferedPipeline(buffer int) *Pipeline {
	sub := reactive.NewBufferSubject(buffer)
	return &Pipeline{
		subject:   sub,
		pipestack: sub,
	}
}
