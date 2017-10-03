package main

import "fmt"

// Workers is a list of Worker
type Workers []Worker

// Worker represents worker
type Worker struct {
	id       int
	jobQueue JobQueue
	quitChan chan bool
}

// NewWorker creates takes a numeric id and a channel w/ worker pool.
func NewWorker(id int, q JobQueue) Worker {
	return Worker{
		id:       id,
		jobQueue: q,
		quitChan: make(chan bool),
	}
}

func (w *Worker) start(s JobQueue) {
	go func() {
		for {
			select {
			case job := <-w.jobQueue:
				result := job.Execute()
				job.Result = result
				s <- job

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
