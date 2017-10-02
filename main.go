package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Job represents Job scheduled in JobQueue
type Job struct {
	Name    string
	Command string
	Delay   int
}

func main() {

	// load config
	confdir := "./conf.d"
	checks, err := loadCfgDir(confdir)
	if err != nil {
		log.Fatalf("Could not load config dir %s: %s", confdir, err.Error())
	}

	// Create the job queue.
	jobQueueLenght := 20
	jobQueue := make(chan Job, jobQueueLenght)

	// run checks
	for _, check := range checks {
		fmt.Println("MAIN: check", check.Name, check.Interval)

		check.Run(jobQueue)
	}

	// TODO run go routine to write results into file
	//resultChannel,resultQuit := newResultWriter()

	var workers []Worker
	maxWorkers := 4
	// Start the dispatcher which will write result string to resultWriter Channel(TODO).
	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(i, jobQueue)
		workers = append(workers, worker)
		worker.start()
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
	shutDown(checks, workers)

	//TODO: waitgroup for worker and checktickers
	time.Sleep(time.Second * 2)
}

func shutDown(c Checks, w Workers) {
	for _, check := range c {
		check.Stop()
	}
	for _, worker := range w {
		worker.stop()
	}
	// resultQuit <- true
}
