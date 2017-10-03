package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

// Sink is resultWriter
type Sink struct {
	Results  JobQueue
	quitChan chan bool
}

// Stop Resultwriter
func (s *Sink) Stop() {
	s.quitChan <- true
}

func startHTTPServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Start Webserver")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func (s *Sink) start() {
	//gaugeVecs := make(map[string]prometheus.GaugeVec)
	gauges := make(map[string]prometheus.Gauge)
	srv := startHTTPServer()

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

// NewSink creates a ResultWriter
func NewSink(c JobQueue) (s Sink) {
	return Sink{Results: c,
		quitChan: make(chan bool),
	}
}
