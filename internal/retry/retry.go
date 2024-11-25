package retry

import (
	"fmt"
	"time"
)

type action func() error
type checkErr func(err error) bool

type retry struct {
	attempt  uint
	max      uint
	action   action
	checkErr checkErr
}

func New(action action, checkErr checkErr, max uint) *retry {
	return &retry{
		attempt:  0,
		checkErr: checkErr,
		max:      max,
		action:   action,
	}
}

func (r *retry) Run() error {
	err := r.action()
	if err != nil && r.checkErr(err) && r.attempt < (r.max-1) {
		r.wait().Run()
		return fmt.Errorf("attempts: %d, %w", (r.attempt + 1), err)
	}

	return err
}

func (r *retry) wait() *retry {
	time.Sleep(r.increment())
	return r
}

func (r *retry) increment() time.Duration {
	d := (1 + 2*time.Duration(r.attempt)) * time.Second
	r.attempt++
	return d
}
