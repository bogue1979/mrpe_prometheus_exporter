package main

import "testing"

func TestWorkerStart(t *testing.T) {
	c := NewBufferedJobQueue(10)
	w := NewWorker(0, c)
	defer w.stop()

	if w.id != 0 {
		t.Errorf("Worker ID to be 0 , got %d\n", w.id)
	}
	c <- Job{Name: "testjob", Command: "testcommand"}
	w.start(c)
}

func TestWorkerLifecycle(t *testing.T) {
	c := NewBufferedJobQueue(10)
	w := NewWorker(0, c)

	w.start(c)
	w.stop()
}
