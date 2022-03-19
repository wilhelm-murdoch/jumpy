package processor

type Results[T any] struct {
	jobs []*Job[T]
}

func (r *Results[T]) Messages() []error {
	var out []error

	r.Each(func(i int, j *Job[T]) {
		out = append(out, j.Messages...)
	})

	return out
}

func (r *Results[T]) Failed() Results[T] {
	return r.Filter(func(j *Job[T]) bool {
		return j.Status == JOB_FAILED
	})
}

func (r *Results[T]) Retried() Results[T] {
	return r.Filter(func(j *Job[T]) bool {
		return j.Status == JOB_RETRIED
	})
}

// Filter returns a new Results struct with items that have passed predicate
// check.
func (r *Results[T]) Filter(f func(*Job[T]) bool) Results[T] {
	var out Results[T]

	r.Each(func(i int, j *Job[T]) {
		if f(j) {
			out.Push(j)
		}
	})

	return out
}

// Push method appends one or more items to the end of an array, returning the
// new length.
func (r *Results[T]) Push(j ...*Job[T]) *Results[T] {
	r.jobs = append(r.jobs, j...)

	return r
}

// Pop removes the last item from an array of Jobs from Results and then returns
// that item.
func (r *Results[T]) Pop() (*Job[T], bool) {
	var job *Job[T]

	if len(r.jobs) == 0 {
		return nil, false
	}

	job = r.jobs[r.Length()-1]
	r.jobs = r.jobs[0 : r.Length()-1]

	return job, true
}

// Length returns the size of the associated list of Jobs.
func (r *Results[T]) Length() int {
	return len(r.jobs)
}

// Map method creates to a new array by using callback invocation result on
// each array item. This method creates a new Results, without mutating the
// original one.
func (r *Results[T]) Map(f func(int, *Job[T]) *Job[T]) Results[T] {
	var out Results[T]

	for i, j := range r.jobs {
		out.Push(f(i, j))
	}

	return out
}

// Each iterates through the associated list of Jobs and executes the specified
// function on each item. This method returns the current instance of Results.
func (r *Results[T]) Each(f func(int, *Job[T])) *Results[T] {
	for i, j := range r.jobs {
		f(i, j)
	}

	return r
}

// Concat merges two slices of Jobs. This method returns the current instance
// Results.
func (r *Results[T]) Concat(jobs []*Job[T]) *Results[T] {
	r.jobs = append(r.jobs, jobs...)

	return r
}
