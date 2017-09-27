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
	Delay   time.Duration
}

// Comment represents comment for prometheus page
type Comment struct {
	Comment string
	Help    string
	Type    string
}

// Check represents mrpe check definition
type Check struct {
	Comment
	Name     string
	Command  string
	Interval time.Duration
	quitChan chan bool
}

// NewCheck to create new Check
func NewCheck() Check {
	return Check{
		Interval: time.Duration(5),
		quitChan: make(chan bool),
	}
}

// Run Check in background
func (c Check) Run(s chan Job) {
	ticker := time.NewTicker(time.Second * c.Interval)
	go func() {
		fmt.Printf("Start %s with Checkinterval %s Seconds\n", c.Name, c.Interval)
		for {
			select {
			case <-ticker.C:
				s <- Job{Command: c.Command, Name: c.Name}
			case <-c.quitChan:
				fmt.Printf("Stopping %s\n", c.Name)
				ticker.Stop()
				return
			}
		}
	}()
}

// Valid checks if required fields set
func (c *Check) Valid() bool {

	if c.Name == "" || c.Command == "" {
		return false
	}
	return true
}

// Checks respresents checklist
type Checks []Check

// Worker represents worker
type Worker struct {
	id       int
	jobQueue chan Job
	quitChan chan bool
}

// NewWorker creates takes a numeric id and a channel w/ worker pool.
func NewWorker(id int, jobQueue chan Job) Worker {
	return Worker{
		id:       id,
		jobQueue: jobQueue,
		quitChan: make(chan bool),
	}
}

func (w Worker) start() {
	go func() {
		for {
			select {
			case job := <-w.jobQueue:
				fmt.Printf("worker%d: started %s Command: %s\n", w.id, job.Name, job.Command)
				//TODO run job.Command and send output to (not yet implemented ) return channel

			case <-w.quitChan:
				fmt.Println("Worker stop:", w.id)
				return
			}
		}
	}()
}

// Workers is a list of Worker
type Workers []Worker

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
		check.quitChan <- true
	}
	for _, worker := range w {
		worker.quitChan <- true
	}
	// resultQuit <- true
}
