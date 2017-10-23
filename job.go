package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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

				gaugeName := fmt.Sprintf("cmk_%s", job.Name)
				_, ok := gauges[gaugeName]
				if !ok {
					gaugeHelp := fmt.Sprintf("Check_MK metrics for %s", job.Name)
					gauges[gaugeName] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: gaugeName,
							Help: gaugeHelp,
						},
						[]string{s.StageKey, "metric"},
					)
					prometheus.MustRegister(gauges[gaugeName])
				}

				if len(job.Perf) > 0 {
					for m, v := range job.Perf {
						gauges[gaugeName].With(
							prometheus.Labels{s.StageKey: s.StageValue, "metric": m}).Set(v)
					}
				}
				gauges[gaugeName].With(
					prometheus.Labels{s.StageKey: s.StageValue, "metric": "exit"}).Set(float64(job.Result.ExitCode))

			case <-s.quitChan:
				log.Info("Stop Webserver")
				if err := srv.Shutdown(nil); err != nil {
					panic(err) // failure/timeout shutting down the server gracefully
				}
				return
			}
		}
	}()
}
