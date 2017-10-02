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
	Stdout   string
	Stderr   string
	Perf     map[string]float64
	Error    error
}

//TODO: run Job and send Result into channel
//      on other side of channel use Result to generate Prometheus Metrics
func runCommand(cmd string, i int) Result {
	var out, eee bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i)*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "sh", "-c", cmd)
	//fmt.Printf("Command: sh -c %s\n", cmd)

	command.Stdout, command.Stderr = &out, &eee
	if err := command.Run(); err != nil {
		return NewResult(-1, "", err.Error(), err)
	}
	stdout, stderr := out.String(), eee.String()

	//fmt.Println("ProcessStateString:", command.ProcessState.String())
	exit, errr := strconv.Atoi(strings.Fields(command.ProcessState.String())[2])
	if errr != nil {
		return NewResult(exit, stdout, stderr, fmt.Errorf("error converting Exitcode %s", errr))
	}

	return NewResult(exit, stdout, stderr, nil)
}

// NewResult creates new Result
func NewResult(exitcode int, stdout, stderr string, err error) Result {
	return Result{ExitCode: exitcode,
		Stdout: stdout,
		Stderr: stderr,
		Error:  err,
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
	return true

}
