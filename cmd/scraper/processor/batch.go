package processor

import (
	"fmt"
	"sync"
)

// A Batch represents a collection of struct Job.
type Batch[T any] struct {
	retryChannel    chan *Job[T]
	finishedChannel chan bool
	Jobs, RetryPool []*Job[T]
}

// Returns a new Batch containing a collection of struct Job.
func NewBatch[T any](jobs []T) Batch[T] {
	batch := Batch[T]{
		retryChannel:    make(chan *Job[T]),
		finishedChannel: make(chan bool),
	}

	for _, job := range jobs {
		batch.Push(NewJob(job))
	}

	return batch
}

// A private method used to monitor the batch processor's signal channels. Will
// collect errors and close all channels processing has completed.
func (b *Batch[T]) messenger() {
	for {
		select {
		case job := <-b.retryChannel:
			b.RetryPool = append(b.RetryPool, job)
		case <-b.finishedChannel:
			close(b.finishedChannel)
			close(b.retryChannel)
			return
		}
	}
}

// Returns the number of jobs associated with the current Batch.
func (b *Batch[T]) Length() int {
	return len(b.Jobs)
}

// Push appends a new Job to the current Batch.
func (b *Batch[T]) Push(job *Job[T]) {
	b.Jobs = append(b.Jobs, job)
}

// Pop attempts to return the last Job in the current Batch
func (b *Batch[T]) Pop() (*Job[T], bool) {
	var job *Job[T]

	if len(b.Jobs) == 0 {
		return nil, false
	}

	job = b.Jobs[b.Length()-1]
	b.Jobs = b.Jobs[0 : b.Length()-1]

	return job, true
}

func (b *Batch[T]) AttemptJob(wg *sync.WaitGroup, p *Processor[T], job *Job[T], worker WorkerHandler[T]) {
	defer wg.Done()

	var err error

	job, err = worker(job)

	fmt.Println("Retries Left:   ", job.RetriesLeft(p.MaxRetries))
	fmt.Println("Retry Attempts: ", job.RetryCount)

	if err != nil {
		job.PushMessage(err)

		if p.CanRetry {
			b.retryChannel <- job
			return
		}
	}

	p.resultChannel <- job
}

// Blocks while executing the specified worker function. Will execute a goroutine
// for each Job while sending any returned errors to the parent Processor's
// errorChannel for handling.
func (b *Batch[T]) Execute(p *Processor[T], worker WorkerHandler[T]) {
	go b.messenger()

	var wg sync.WaitGroup

	wg.Add(b.Length())

	for _, job := range b.Jobs {
		go b.AttemptJob(&wg, p, job, worker)
	}

	wg.Wait()

	b.finishedChannel <- true
}
