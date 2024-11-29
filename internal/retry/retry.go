package retry

import (
	"fmt"
	"time"
)

type Action func() error
type CheckErr func(err error) bool

type retry struct {
	action   Action
	checkErr CheckErr
	attempt  uint
	max      uint
}

func New(action Action, checkErr CheckErr, count uint) *retry {
	return &retry{
		attempt:  1,
		checkErr: checkErr,
		max:      count,
		action:   action,
	}
}

func (r *retry) Run() error {
	var err error
	for {
		err = r.action()
		if err == nil || !r.checkErr(err) || r.attempt >= r.max {
			break
		}

		r.wait()
		r.attempt++
	}

	if err != nil {
		return fmt.Errorf("attempts: %d, %w", r.attempt, err)
	}

	return nil
}

func (r *retry) wait() {
	time.Sleep(r.duration())
}

func (r *retry) duration() time.Duration {
	return (1 + 2*time.Duration(r.attempt-1)) * time.Second
}
