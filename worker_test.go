package main

import "testing"

func TestWorkerStart(t *testing.T) {
	c := NewBufferedJobQueue(10)
	w := NewWorker(0, c)
	defer w.stop()

	if w.id != 0 {
		t.Errorf("expected worker id to be 0 , got %d\n", w.id)
	}
	c <- Job{Name: "testjob", Command: "testcommand"}
	w.start(c)
}

func TestWorkerLifecycle(t *testing.T) {
	c := NewBufferedJobQueue(10)
	w := NewWorker(0, c)

	if w.running {
		t.Errorf("expected worker running to be false. got: %v", w.running)
	}

	w.start(c)
	if !w.running {
		t.Errorf("expected worker running to be true. got: %v", w.running)
	}
	w.stop()
}
