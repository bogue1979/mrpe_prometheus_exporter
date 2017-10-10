package main

import "testing"

func TestWorkerStart(t *testing.T) {
	c := newBufferedJobQueue(10)
	w := newWorker(0, c)
	defer w.stop()

	if w.id != 0 {
		t.Errorf("expected worker id to be 0 , got %d\n", w.id)
	}
	c <- Job{Name: "testjob", Command: "echo one"}
	w.start(c)
}

func TestWorkerLifecycle(t *testing.T) {
	c := newBufferedJobQueue(10)
	w := newWorker(0, c)

	if w.running {
		t.Errorf("expected worker running to be false. got: %v", w.running)
	}

	w.start(c)
	if !w.running {
		t.Errorf("expected worker running to be true. got: %v", w.running)
	}
	w.stop()
}

func TestWorkerExecute(t *testing.T) {
	c := newBufferedJobQueue(10)
	w := newWorker(0, c)
	defer w.stop()

	c <- Job{Name: "testjob", Command: "echo one"}
	w.start(c)
	j := <-c

	if j.ExitCode != 0 {
		t.Errorf("expected to execute with exitcode of 0. got %d", j.ExitCode)
	}
}
