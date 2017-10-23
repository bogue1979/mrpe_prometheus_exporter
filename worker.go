package main

import log "github.com/sirupsen/logrus"

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
				if job.Result.Error != nil {
					log.Debugf("error in result for %s: %s", job.Name, job.Result.Error)
				}
				job.Result = result
				if err := job.PerformanceData(); err != nil {
					log.Errorf("problem getting PerformanceData: %s", err)
				}
				s <- job

			case <-w.quitChan:
				log.Infof("Worker stop: %d", w.id)
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	w.running = false
	w.quitChan <- true
}
