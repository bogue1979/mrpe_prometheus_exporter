package main

import "fmt"

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

func (w *Worker) start() {
	go func() {
		for {
			select {
			case job := <-w.jobQueue:
				result := runCommand(job.Command, job.Delay)
				if result.PerformanceData() {
					//fmt.Println("generated PerformanceData")
				}
				//fmt.Printf("result: %#v", result)
				//TODO run job.Command and send output to (not yet implemented ) return channel

			case <-w.quitChan:
				fmt.Println("Worker stop:", w.id)
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	w.quitChan <- true
}

// Workers is a list of Worker
type Workers []Worker
