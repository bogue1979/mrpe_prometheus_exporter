package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// JobQueue represents as named
type JobQueue chan Job

// NewJobQueue creates a new JobQueue
func NewJobQueue() JobQueue {
	return make(chan Job)
}

// NewBufferedJobQueue returns buffered JobQueue
func NewBufferedJobQueue(i int) JobQueue {
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
	Results  JobQueue
	quitChan chan bool
	running  bool
}

// Stop Resultwriter
func (s *resultWriter) Stop() {
	s.running = false
	s.quitChan <- true
}

func (s *resultWriter) start() {
	//gaugeVecs := make(map[string]prometheus.GaugeVec)
	gauges := make(map[string]prometheus.Gauge)
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
				g, ok := gauges[gaugeExitCodeName]
				if !ok {
					//TODO use help string in job.Comment
					gaugeHelp := fmt.Sprintf("Check exitcode for %s", job.Name)
					gauges[gaugeExitCodeName] = prometheus.NewGauge(prometheus.GaugeOpts{
						Name: gaugeExitCodeName,
						Help: gaugeHelp,
					})
					prometheus.MustRegister(gauges[gaugeExitCodeName])
					gauges[gaugeExitCodeName].Set(float64(job.Result.ExitCode))
				} else {
					g.Set(float64(job.Result.ExitCode))
				}

				if job.Duration > 0 {
					gaugeDurationName := fmt.Sprintf("cmk_%s_duration_ns", job.Name)
					g, ok := gauges[gaugeDurationName]
					if !ok {
						gaugeHelp := fmt.Sprintf("Runtime in ns for %s", job.Name)
						gauges[gaugeDurationName] = prometheus.NewGauge(prometheus.GaugeOpts{
							Name: gaugeDurationName,
							Help: gaugeHelp,
						})
						prometheus.MustRegister(gauges[gaugeDurationName])
						gauges[gaugeExitCodeName].Set(float64(job.Result.Duration))
					} else {
						g.Set(float64(job.Result.Duration))
					}

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
