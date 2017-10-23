package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "undef"

func setLogLevel(s string) {
	switch s {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.Errorf("cannot use %s as log.lvl, will use default info level", s)
	}
}

func main() {

	var stagekey = flag.String("env.key", "stage", "environment differentiator")
	var stageval = flag.String("env.val", "dev", "environment name")
	var confdir = flag.String("conf.dir", "./conf.d", "directory with mrpe config files")
	var versionstring = flag.Bool("version", false, "show version and exit")
	var loglevel = flag.String("log.lvl", "info", "loglevel from [debug,info,warn,error,fatal,panic]")
	var logjson = flag.Bool("log.json", false, "log as json")

	flag.Parse()

	setLogLevel(*loglevel)
	if *logjson {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if *versionstring {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	checks, err := loadCfgDir(*confdir)
	if err != nil {
		log.Fatalf("Could not load config dir %s: %s", *confdir, err.Error())
	}

	// Create the job queue.
	jobQueueLenght := 20
	jobQueue := newBufferedJobQueue(jobQueueLenght)

	// run checks
	for _, check := range checks {
		check.start(jobQueue)
	}

	// resultWriter
	resultChan := newJobQueue()
	sink := newResultWriter(resultChan, *stagekey, *stageval)
	sink.start()

	var workers []Worker
	maxWorkers := 4
	// Start the dispatcher which will write result string to resultWriter Channel
	for i := 0; i < maxWorkers; i++ {
		worker := newWorker(i, jobQueue)
		workers = append(workers, worker)
		worker.start(sink.Results)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Infof("Got %s signal", sig)
		done <- true
	}()

	log.Info("Program started")

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
	log.Info("Start Webserver")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}
