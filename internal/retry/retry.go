package retry

import (
	"fmt"
	"time"
)

type action func() error
type checkErr func(err error) bool

type retry struct {
	action   action
	checkErr checkErr
	attempt  uint
	max      uint
}

func New(action action, checkErr checkErr, count uint) *retry {
	return &retry{
		attempt:  1,
		checkErr: checkErr,
		max:      count,
		action:   action,
	}
}

func (r *retry) Run() error {
	var err error
	for r.attempt <= r.max {
		err = r.action()
		if err == nil || !r.checkErr(err) || r.attempt == r.max {
			break
		}

		r.wait()
	}

	if err != nil {
		return fmt.Errorf("attempts: %d, %w", r.attempt, err)
	}

	return nil
}

func (r *retry) wait() {
	time.Sleep(r.increment())
}

func (r *retry) increment() time.Duration {
	d := (1 + 2*time.Duration(r.attempt-1)) * time.Second
	r.attempt++
	return d
}
