package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
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

func sanitizeLabel(label, unit string) string {
	var suffix string

	switch unit {
	case "%":
		suffix = "_percent"
	case "":
		suffix = ""
	default:
		suffix = "_" + unit
	}
	return strings.ToLower(strings.Replace(label, " ", "_", -1) + suffix)
}

func perfstringMap(perflist []string) (r map[string]float64, err error) {
	r = make(map[string]float64)
	// ["match","label","value","unit","warn","crit","min","max"]
	value, err := strconv.ParseFloat(perflist[2], 64)
	if err != nil {
		return r, fmt.Errorf("could not convert value for perfdata %s: %s", perflist[1], err)
	}
	r[sanitizeLabel(perflist[1], perflist[3])] = value

	index := map[int]string{4: "warn",
		5: "crit",
		6: "min",
		7: "max"}

	for i := range perflist[4:] {
		if perflist[i+4] != "" {
			val, err := strconv.ParseFloat(perflist[i+4], 64)
			if err != nil {
				return r, fmt.Errorf("could not convert %s for perfdata %s: %s", index[i+4], perflist[1], err)
			}
			r[sanitizeLabel(perflist[1], index[i+4])] = val
		}
	}
	return r, nil
}

// PerformanceData parses Stdout String to Perf map
func (r *Result) PerformanceData() (err error) {
	perfsplit := strings.Split(r.Stdout, "|")
	if len(perfsplit) != 2 {
		return nil
	}
	r.Perf = make(map[string]float64)
	if r.Duration > 0 {
		r.Perf["duration"] = float64(r.Duration)
	}

	s := strings.Split(perfsplit[1], " ")
	for _, i := range s {

		rx := regexp.MustCompile("^'?(.*)'?=([1-9.,]+)([a-zA-Z%]+)*;?([0-9.,]*);?([0-9.,]*);?([0-9.,]*);?([0-9.,]*)")
		perflist := rx.FindStringSubmatch(i)
		if len(perflist) > 0 {
			fmt.Printf("%q\n", perflist)
			perfmap, err := perfstringMap(perflist)
			if err != nil {
				return fmt.Errorf("Problems generating PerformanceData: %s", err)
			}
			for k, v := range perfmap {
				r.Perf[k] = v
			}
		}
	}
	return nil
}
