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
