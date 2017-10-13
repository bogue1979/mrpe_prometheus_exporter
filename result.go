package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Result represents CheckResult
type Result struct {
	ExitCode int
	Duration int64
	Stdout   string
	Stderr   string
	Perf     map[string]float64
	Error    error
}

func runCommand(cmd string, i int64) Result {
	var out, eee bytes.Buffer
	var result Result

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i)*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "/bin/sh", "-c", cmd)

	t1 := time.Now()
	command.Stdout, command.Stderr = &out, &eee
	t2 := time.Now()
	if err := command.Run(); err != nil {
		return NewResult(-1, "", err.Error(), err)
	}
	stdout, stderr := out.String(), eee.String()

	exit, errr := strconv.Atoi(strings.Fields(command.ProcessState.String())[2])
	if errr != nil {
		return NewResult(exit, stdout, stderr, fmt.Errorf("error converting Exitcode %s", errr))
	}
	result = NewResult(exit, stdout, stderr, nil)
	result.Duration = int64(t2.Sub(t1))

	return result
}

// NewResult creates new Result
func NewResult(exitcode int, stdout, stderr string, err error) Result {
	return Result{ExitCode: exitcode,
		Stdout: stdout,
		Stderr: stderr,
		Error:  err,
	}
}

// newResultWriter creates a ResultWriter
func newResultWriter(c JobQueue, sKey, sValue string) (s resultWriter) {
	return resultWriter{Results: c,
		quitChan:   make(chan bool),
		StageKey:   sKey,
		StageValue: sValue,
	}
}

// PerformanceData parses Stdout String to Perf map
func (r *Result) PerformanceData() (ok bool) {
	perfsplit := strings.Split(r.Stdout, "|")
	if len(perfsplit) != 2 {
		return false
	}
	r.Perf = make(map[string]float64)
	s := strings.Split(perfsplit[1], ",")
	for _, i := range s {
		labelValue := strings.Split(strings.TrimSpace(i), "=")
		if len(labelValue) != 2 {
			return false
		}
		f, err := strconv.ParseFloat(labelValue[1], 64)
		if err != nil {
			return false
		}
		r.Perf[labelValue[0]] = f
	}
	if r.Duration > 0 {
		r.Perf["duration"] = float64(r.Duration)
	}
	return true
}
