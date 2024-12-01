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
		attempt:  0,
		checkErr: checkErr,
		max:      count,
		action:   action,
	}
}

func (r *retry) Run() error {
	var err error
	for ; r.attempt < r.max; r.attempt++ {
		r.wait()

		err = r.action()
		if err == nil || !r.checkErr(err) {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("attempts: %d, %w", r.attempt, err)
	}

	return nil
}

func (r *retry) wait() {
	d := r.getDuration()
	if d > 0 {
		time.Sleep(d)
	}
}

func (r *retry) getDuration() time.Duration {
	return (1 + 2*time.Duration(r.attempt-1)) * time.Second
}
