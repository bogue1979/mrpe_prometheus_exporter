package main

import (
	"testing"
)

func TestCheckRunner(t *testing.T) {
	jobs := newJobQueue()
	c := NewCheck()
	c.Name = "mycheck"
	c.Command = "mycheck command"
	c.Interval = 1

	c.start(jobs)

	j := <-jobs
	if j.Name != "mycheck" {
		t.Errorf("Expect Jobname testcheck , got %s", j.Name)
	}
	c.Stop()
}
