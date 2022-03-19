package processor

import (
	"math"
	"time"
)

type ProcessorOptionals struct {
	CanRetry        bool
	MaxRetries      int
	ThrottleRetries bool
	RetryInterval   time.Duration
	BatchSize       int
}

type Option func(*ProcessorOptionals)

func CanRetry(canRetry bool) Option {
	return func(args *ProcessorOptionals) {
		args.CanRetry = canRetry
	}
}

func MaxRetries(maxRetries int) Option {
	return func(args *ProcessorOptionals) {
		args.MaxRetries = maxRetries
	}
}

func ThrottleRetries(throttleRetries bool) Option {
	return func(args *ProcessorOptionals) {
		args.ThrottleRetries = throttleRetries
	}
}

func RetryInterval(retryInterval time.Duration) Option {
	return func(args *ProcessorOptionals) {
		args.RetryInterval = retryInterval
	}
}

func BatchSize(batchSize int) Option {
	return func(args *ProcessorOptionals) {
		args.BatchSize = batchSize
	}
}

type WorkerHandler[T any] func(*Job[T]) (*Job[T], error)

// Represents a batch processor.
type Processor[T any] struct {
	CanRetry        bool
	MaxRetries      int
	ThrottleRetries bool
	RetryInterval   time.Duration
	Errors          []error
	Results         []*Job[T]
	resultChannel   chan *Job[T]
	errorChannel    chan error
	finishedChannel chan bool
	BatchSize       int
	BatchCount      int
	Batches         []Batch[T]
}

// A private method used to monitor the batch processor's signal channels. Will
// collect errors and close all channels processing has completed.
func (p *Processor[T]) messenger() {
	for {
		select {
		case result := <-p.resultChannel:
			p.Results = append(p.Results, result)
		case err := <-p.errorChannel:
			p.Errors = append(p.Errors, err)
		case <-p.finishedChannel:
			close(p.finishedChannel)
			close(p.errorChannel)
			close(p.resultChannel)
			return
		}
	}
}

// Will begin batch processing by iterating through the current collection of
// batch jobs which will ultimately block and execute the specified worker
// function.
func (p *Processor[T]) Execute(worker WorkerHandler[T]) []*Job[T] {
	go p.messenger()

	for _, b := range p.Batches {
		b.Execute(p, worker)
	}

	p.finishedChannel <- true

	return p.Results
}

// Returns a new Processor whose state will be primed for executing batch jobs.
// Job distribution into a collection of Batches are calculated here before
// returning.
func NewProcessor[T any](jobs []T, options ...Option) Processor[T] {
	defaults := &ProcessorOptionals{
		CanRetry:        true,
		MaxRetries:      3,
		ThrottleRetries: false,
		RetryInterval:   time.Second * 5,
		BatchSize:       10,
	}

	for _, o := range options {
		o(defaults)
	}

	jobsCount := len(jobs)

	if defaults.BatchSize > jobsCount {
		defaults.BatchSize = jobsCount
	}

	batchCount := int(math.Ceil(float64(len(jobs)) / float64(defaults.BatchSize)))

	var batches []Batch[T]

	offset, limit := 0, defaults.BatchSize
	for i := 0; i < batchCount; i++ {
		batches = append(batches, NewBatch(jobs[offset:limit]))
		offset, limit = limit, limit+defaults.BatchSize

		if limit > jobsCount {
			limit = jobsCount
		}
	}

	return Processor[T]{
		CanRetry:        defaults.CanRetry,
		MaxRetries:      defaults.MaxRetries,
		ThrottleRetries: defaults.ThrottleRetries,
		RetryInterval:   defaults.RetryInterval,
		BatchSize:       defaults.BatchSize,
		BatchCount:      batchCount,
		Batches:         batches,
		errorChannel:    make(chan error),
		finishedChannel: make(chan bool),
		resultChannel:   make(chan *Job[T]),
	}
}
