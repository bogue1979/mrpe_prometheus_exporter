package main

import (
	"fmt"
	"time"
)

// Comment represents comment for prometheus page
type Comment struct {
	Comment string
	Help    string
	Type    string
}

// Check represents mrpe check definition
type Check struct {
	Comment
	Name     string
	Command  string
	Interval int64
	quitChan chan bool
}

// NewCheck to create new Check
func NewCheck() Check {
	return Check{
		Interval: 5,
		quitChan: make(chan bool),
	}
}

// Stop check routine
func (c *Check) Stop() {
	c.quitChan <- true
}

// Run Check in background
func (c Check) Run(s JobQueue) {

	ticker := time.NewTicker(time.Second * time.Duration(c.Interval))
	go func() {
		fmt.Printf("Start %s with Checkinterval %d Seconds\n", c.Name, c.Interval)
		for {
			select {
			case <-ticker.C:
				s <- Job{Command: c.Command, Name: c.Name}
			case <-c.quitChan:
				fmt.Printf("Stopping %s\n", c.Name)
				ticker.Stop()
				return
			}
		}
	}()
	return
}

// Valid checks if required fields set
func (c *Check) Valid() bool {

	if c.Name == "" || c.Command == "" {
		return false
	}
	return true
}

// Checks respresents checklist
type Checks []Check
