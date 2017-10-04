package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	//TODO flag for environment vars
	stagekey := "e42stage"
	stageval := "dev"

	//TODO flag for confdir
	confdir := "./conf.d"
	checks, err := loadCfgDir(confdir)
	if err != nil {
		log.Fatalf("Could not load config dir %s: %s", confdir, err.Error())
	}

	// Create the job queue.
	jobQueueLenght := 20
	jobQueue := NewBufferedJobQueue(jobQueueLenght)

	// run checks
	for _, check := range checks {
		check.start(jobQueue)
	}

	// resultWriter
	resultChan := NewJobQueue()
	sink := newResultWriter(resultChan, stagekey, stageval)
	sink.start()

	var workers []Worker
	maxWorkers := 4
	// Start the dispatcher which will write result string to resultWriter Channel
	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(i, jobQueue)
		workers = append(workers, worker)
		worker.start(sink.Results)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("Got %s signal\n", sig)
		done <- true
	}()

	fmt.Println("Program started")

	// wait for Os.Signal to shutown service
	<-done
	shutDown(checks, workers, sink)

	//TODO: waitgroup for worker and checktickers
	time.Sleep(time.Second * 2)
}

func shutDown(c Checks, w Workers, s resultWriter) {
	for _, check := range c {
		check.Stop()
	}
	for _, worker := range w {
		worker.stop()
	}
	s.Stop()
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
