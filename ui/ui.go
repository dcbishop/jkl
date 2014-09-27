package ui

import "github.com/dcbishop/jkl/service"

// UI handles input and displaying the app.
type UI interface {
	service.Service
	Input
}

// Input implementers provide a chanel that generates events.
type Input interface {
	Events() <-chan Event
}

// Event holds information about an event
type Event struct {
	Data interface{}
}
