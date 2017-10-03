package main

import "testing"

func TestNewSink(t *testing.T) {
	resultchannel := NewJobQueue()
	sink := newResultWriter(resultchannel)
	sink.start()
	if !sink.running {
		t.Errorf("expected resultWriter.running to be true. got %v", sink.running)
	}
	sink.Stop()
}

func TestExecute(t *testing.T) {
	j := Job{Command: "exit 0", Delay: 1}
	r := j.Execute()
	if r.ExitCode != 0 {
		t.Errorf("expected exitcode to be 0. got: %v", r.ExitCode)
	}
}
