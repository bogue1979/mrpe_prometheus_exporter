package main

import (
	"testing"
)

func TestCommand(t *testing.T) {
	gotresult := runCommand("exit 0", 2)

	if gotresult.Error != nil {
		t.Errorf("Expectet no error, got %s", gotresult.Error)
	}
}

func TestCommandPerfdataMegaByte(t *testing.T) {
	expPerf := map[string]float64{
		"foo_mb":   1.3,
		"foo_warn": 1.4,
		"foo_crit": 1.5,
		"foo_min":  0.5,
		"foo_max":  10,
	}
	gotresult := runCommand("echo 'testoutput | foo=1.3MB;1.4;1.5;0.5;10'", 2)
	gotresult.PerformanceData()

	for k := range expPerf {
		_, ok := gotresult.Perf[k]
		if !ok {
			t.Errorf("Expected to have %s in PerformanceData map, but it is not", k)
		}
	}
}

func TestCommandPerfdataPercent(t *testing.T) {
	expPerf := map[string]float64{
		"foo_percent": 1.3,
		"bar":         1,
		"duration":    5,
	}
	gotresult := runCommand("echo 'testoutput | foo=1.3% bar=1'", 2)
	gotresult.PerformanceData()

	for k := range expPerf {
		_, ok := gotresult.Perf[k]
		if !ok {
			t.Errorf("Expected to have %s in PerformanceData map, but it is not", k)
		}
	}
}

func TestNoPerformanceData(t *testing.T) {
	gotresult := runCommand("echo 'testoutput ", 2)
	err := gotresult.PerformanceData()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestCommandNotFound(t *testing.T) {
	gotresult := runCommand("testcommand not found", 1)

	if gotresult.ExitCode == 0 {
		t.Errorf("Expected results Exitcode not to be 0 but is was %v", gotresult.ExitCode)
	}
}

//func TestCommandWrongPerfdata(t *testing.T) {
//	cmd2 := runCommand("echo 'foo | bar=a'", 2)
//	cmd3 := runCommand("echo 'foo | bar=baz'", 2)
//
//	for _, cmd := range []Result{cmd2, cmd3} {
//
//		err := cmd.PerformanceData()
//		if err == nil {
//			t.Errorf("Expected error in PerformanceData() for command %#v", cmd)
//		}
//	}
//}

func TestCommandTimeout(t *testing.T) {

	gotresult := runCommand("sleep 5", 2)

	if gotresult.ExitCode != -1 {
		t.Errorf("Expected results ExitCode to be -1 but is was %d", gotresult.ExitCode)
	}
}
