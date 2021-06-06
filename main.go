package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	types   = []string{"emai", "deactivation", "activation", "transaction", "customer_renew", "order_processed"}
	workers = 0

	totalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "worker",
			Subsystem: "jobs",
			Name:      "processed_total",
			Help:      "Total number of jobs processed by the workers",
		},
		[]string{"worker_id", "type"},
	)

	inflightCounterVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "worker",
			Subsystem: "jobs",
			Name:      "inflight",
			Help:      "Number of jobs inflight",
		},
		[]string{"type"},
	)

	processingTimeVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "worker",
			Subsystem: "jobs",
			Name:      "process_time_seconds",
			Help:      "Amount of time spent processing jobs",
		},
		[]string{"worker_id", "type"},
	)
)

func init() {
	flag.IntVar(&workers, "workers", 10, "Number of workers to use")
}

func getType() string {
	return types[rand.Int()%len(types)]
}

func main() {
	flag.Parse()

	prometheus.MustRegister(
		totalCounterVec,
		inflightCounterVec,
		processingTimeVec,
	)

	// create a channel with a 10,000 Job buffer
	jobsChannel := make(chan *Job, 10000)

	go startJobProcessor(jobsChannel)

	go createJobs(jobsChannel)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("[INFO] starting HTTP server on port :9009")
	log.Fatal(http.ListenAndServe(":9009", nil))
}

type Job struct {
	Type  string
	Sleep time.Duration
}

// makeJob creates a new job with a random sleep time between 10 ms and 4000ms.
func makeJob() *Job {
	return &Job{
		Type:  getType(),
		Sleep: time.Duration(rand.Int()%100+10) * time.Millisecond,
	}
}

func startJobProcessor(jobs <-chan *Job) {
	log.Printf("[INFO] starting %d workers\n", workers)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			// Start the worker
			startWorker(workerID, jobs)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func createJobs(jobs chan<- *Job) {
	for {
		// Create random job.
		job := makeJob()
		// track the job in the inflight tracker
		inflightCounterVec.WithLabelValues(job.Type).Inc()
		// Send the job down the channel
		jobs <- job
		// Don't pile up too quickly
		time.Sleep(5 * time.Millisecond)
	}
}

// creates a worker that pulls jobs from the job channel
func startWorker(workerID int, jobs <-chan *Job) {
	for job := range jobs {
		startTime := time.Now()

		// fake processing the request
		time.Sleep(job.Sleep)
		log.Printf("[%d][%s] Processed job in %0.3f seconds", workerID, job.Type, time.Since(startTime).Seconds())
		totalCounterVec.WithLabelValues(strconv.FormatInt(int64(workerID), 10), job.Type).Inc()
		// Decrement the inflight counter.
		inflightCounterVec.WithLabelValues(job.Type).Dec()

		processingTimeVec.WithLabelValues(strconv.FormatInt(int64(workerID), 10), job.Type).Observe(time.Since(startTime).Seconds())
	}
}
