package service

import (
	"errors"
	"sync/atomic"
	"time"
)

// Service implementations have a Run() loop that does stuff and blocks untill Stop() is called
type Service interface {
	Run()
	Stop()
	Running() bool
}

// WaitUntilRunning will wait until service.Running() returns true,
// will return an error if it took longer than timeout.
func WaitUntilRunning(service Service, timeout time.Duration) error {
	return WaitUntil(service, true, timeout)
}

// WaitUntilStopped will wait until service.Running() returns false,
// will return an error if it took longer than timeout.
func WaitUntilStopped(service Service, timeout time.Duration) error {
	return WaitUntil(service, false, timeout)
}

// WaitUntil will wait until service.Running() returns state,
// will return an error if it took longer than timeout.
func WaitUntil(service Service, state bool, timeout time.Duration) error {
	start := time.Now()

	for !service.Running() == state {
		time.Sleep(time.Millisecond)
		if time.Since(start) >= timeout {
			return errors.New("Time out")
		}
	}

	return nil
}

// State tracks the running state for Service implementors in a threadsafe way.
type State uint32

// SetRunning sets the running state to true.
func (state *State) SetRunning() error {
	return state.SetState(true)
}

// SetStopped sets the running state to true.
func (state *State) SetStopped() error {
	return state.SetState(false)
}

// SetState sets the running state of the service to the given value.
func (state *State) SetState(newState bool) error {
	newNum := uint32(0)

	if newState {
		newNum = 1
	}

	previousState := atomic.SwapUint32((*uint32)(state), newNum)

	if previousState == 0 && !newState ||
		previousState == 1 && newState {
		return errors.New("Already set.")
	}

	return nil
}

// Running returns true if the state of the service is running.
func (state *State) Running() bool {
	return *(*uint32)(state) == 1
}
