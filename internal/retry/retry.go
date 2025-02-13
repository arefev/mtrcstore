// The retry package repeats the operation until it reaches the number of attempts specified in count.
package retry

import (
	"fmt"
	"time"
)

// The action function that retry calls.
type Action func() error

// A function that checks for an error.
type CheckErr func(err error) bool

type retry struct {
	action   Action // The action function that retry calls.
	checkErr CheckErr // A function that checks for an error.
	attempt  uint // Number of attempts made
	max      uint // Number of maximum attempts
}

// Create new retry object.
func New(action Action, checkErr CheckErr, count uint) *retry {
	return &retry{
		attempt:  1,
		checkErr: checkErr,
		max:      count,
		action:   action,
	}
}

// Run calls Action as long as CheckErr == true and the number of attempts is less than the maximum set.
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

// Stops program execution for the period specified in getDuration.
func (r *retry) wait() {
	time.Sleep(r.getDuration())
}

// Specifies the period for which the program execution will stop.
func (r *retry) getDuration() time.Duration {
	return (1 + 2*time.Duration(r.attempt-1)) * time.Second
}
