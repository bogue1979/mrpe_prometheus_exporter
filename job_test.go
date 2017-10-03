package main

import "testing"

func TestNewSink(t *testing.T) {
	resultchannel := NewJobQueue()
	sink := NewSink(resultchannel)
	sink.start()
	sink.Stop()
}
