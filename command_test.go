package main

import (
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {
	expresult := Result{ExitCode: 0, Stdout: "", Stderr: "", Perf: nil, Error: nil}

	gotresult := runCommand("exit 0", 2)

	if !reflect.DeepEqual(gotresult, expresult) {
		t.Errorf("Expected result tp be %#v but is was %#v", expresult, gotresult)
	}
}

func TestCommandPerfdata(t *testing.T) {
	expPerf := map[string]float64{
		"foo": 1.3,
		"bar": 1,
	}
	gotresult := runCommand("echo 'testoutput | foo=1.3, bar=1'", 2)
	gotresult.PerformanceData()

	if !reflect.DeepEqual(gotresult.Perf, expPerf) {
		t.Errorf("Expected result tp be %#v but is was %#v", expPerf, gotresult.Perf)
	}
}

func TestCommandNotFound(t *testing.T) {
	gotresult := runCommand("testcommand not found", 1)

	if gotresult.ExitCode != -1 {
		t.Errorf("Expected results Exitcode to be -1 but is was %v", gotresult.ExitCode)
	}
}

func TestCommandWrongPerfdata(t *testing.T) {
	cmd1 := runCommand("echo foo bar", 2)
	cmd2 := runCommand("echo 'foo | bar'", 2)
	cmd3 := runCommand("echo 'foo | bar=baz'", 2)

	for _, cmd := range []Result{cmd1, cmd2, cmd3} {

		ok := cmd.PerformanceData()
		if ok {
			t.Errorf("Expected no PerformanceData for command %#v", cmd)
		}
	}
}

func TestCommandTimeout(t *testing.T) {

	gotresult := runCommand("sleep 5", 2)

	if gotresult.ExitCode != -1 {
		t.Errorf("Expected results ExitCode to be -1 but is was %d", gotresult.ExitCode)
	}
}
