package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// load config
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
		check.Run(jobQueue)
	}

	// TODO run go routine to write results into file
	//resultQuit := newResultWriter()
	resultChan := NewJobQueue()
	sink := NewSink(resultChan)
	sink.start()

	var workers []Worker
	maxWorkers := 4
	// Start the dispatcher which will write result string to resultWriter Channel(TODO).
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

func shutDown(c Checks, w Workers, s Sink) {
	for _, check := range c {
		check.Stop()
	}
	for _, worker := range w {
		worker.stop()
	}
	s.Stop()
}
