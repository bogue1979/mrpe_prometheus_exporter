package main

import "testing"

func TestNewSink(t *testing.T) {
	resultchannel := NewJobQueue()
	sink := newResultWriter(resultchannel)
	sink.start()
	sink.Stop()
}
