package main

import "fmt"

// Workers is a list of Worker
type Workers []Worker

// Worker represents worker
type Worker struct {
	id       int
	running  bool
	jobQueue JobQueue
	quitChan chan bool
}

// newWorker creates takes a numeric id and a channel w/ worker pool.
func newWorker(id int, q JobQueue) Worker {
	return Worker{
		id:       id,
		running:  false,
		jobQueue: q,
		quitChan: make(chan bool),
	}
}

func (w *Worker) start(s JobQueue) {
	w.running = true
	go func() {
		for {
			select {
			case job := <-w.jobQueue:
				result := job.Execute()
				job.Result = result
				if err := job.PerformanceData(); err != nil {
					//TODO: Logging
					fmt.Println("problem getting PerformanceData", err)
				}
				s <- job

			case <-w.quitChan:
				fmt.Println("Worker stop:", w.id)
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	w.running = false
	w.quitChan <- true
}
