package processor

const (
	JOB_FAILED = 1 << iota
	JOB_PASSED
	JOB_PROCESSING
	JOB_RETRIED
)

// A Job represents a chunk of data that requires batch processing.
type Job[T any] struct {
	Body       T
	RetryCount int
	Status     int
	Messages   []error
}

func (j *Job[T]) PushMessage(message error) {
	j.Messages = append(j.Messages, message)
}

// RetriesLeft returns the remaining number of available retry attempts for the
// current Job.
func (j *Job[T]) RetriesLeft(max int) int {
	return max - j.RetryCount
}

// Returns a new Job and assigns the supplied generic body value to the associated
// Body field.
func NewJob[T any](body T) *Job[T] {
	return &Job[T]{Body: body}
}
