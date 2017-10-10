package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// JobQueue represents as named
type JobQueue chan Job

// newJobQueue creates a new JobQueue
func newJobQueue() JobQueue {
	return make(chan Job)
}

// newBufferedJobQueue returns buffered JobQueue
func newBufferedJobQueue(i int) JobQueue {
	return make(chan Job, i)
}

// Job represents Job scheduled in JobQueue
type Job struct {
	Comment
	Name    string
	Command string
	Delay   int64
	Result
}

// Execute Job in shell
func (j Job) Execute() (result Result) {
	result = runCommand(j.Command, j.Delay)
	return result
}

// resultWriter is resultWriter
type resultWriter struct {
	Results    JobQueue
	quitChan   chan bool
	running    bool
	StageKey   string
	StageValue string
}

// Stop Resultwriter
func (s *resultWriter) Stop() {
	s.running = false
	s.quitChan <- true
}

func (s *resultWriter) start() {
	//gaugeVecs := make(map[string]prometheus.GaugeVec)
	gauges := make(map[string]*prometheus.GaugeVec)
	srv := startHTTPServer()
	s.running = true

	go func() {
		for {
			select {
			case job := <-s.Results:

				if len(job.Perf) > 0 {
					continue
				}
				// ExitcodeName
				gaugeExitCodeName := fmt.Sprintf("cmk_%s_exit_code", job.Name)
				_, ok := gauges[gaugeExitCodeName]
				if !ok {
					//TODO use help string from job.Comment
					gaugeHelp := fmt.Sprintf("check exitcode for %s", job.Name)
					gauges[gaugeExitCodeName] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: gaugeExitCodeName,
							Help: gaugeHelp,
						},
						[]string{s.StageKey},
					)
					prometheus.MustRegister(gauges[gaugeExitCodeName])
				}
				gauges[gaugeExitCodeName].With(
					prometheus.Labels{s.StageKey: s.StageValue}).Set(float64(job.Result.ExitCode))

				if job.Duration > 0 {
					gaugeDurationName := fmt.Sprintf("cmk_%s_duration_ns", job.Name)
					_, ok := gauges[gaugeDurationName]
					if !ok {
						gaugeHelp := fmt.Sprintf("runtime in ns for %s", job.Name)
						gauges[gaugeDurationName] = prometheus.NewGaugeVec(
							prometheus.GaugeOpts{
								Name: gaugeDurationName,
								Help: gaugeHelp,
							},
							[]string{s.StageKey},
						)
						prometheus.MustRegister(gauges[gaugeDurationName])
					}
					gauges[gaugeDurationName].With(
						prometheus.Labels{s.StageKey: s.StageValue}).Set(float64(job.Result.Duration))
				}
			case <-s.quitChan:
				fmt.Println("Stop Webserver")
				if err := srv.Shutdown(nil); err != nil {
					panic(err) // failure/timeout shutting down the server gracefully
				}
				return
			}
		}
	}()
}
