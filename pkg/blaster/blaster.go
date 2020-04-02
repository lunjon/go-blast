package blaster

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lunjon/go-blast/pkg/util"
)

// Blaster represents a so called blaster, it runs for a given duration
// and each time it's ticker emits a value.
type Blaster struct {
	running bool
	id      string

	config *Configuration
	ticker   *time.Ticker
	httpClient  *http.Client

	// For the report
	totalRequests      int
	successfulRequests int

	done chan bool
	wg   *sync.WaitGroup
}

// NewBlaster creates a new blaster that can be started.
// The duration specifies how long it will run,
// the period how long time between each tick,
// and a sync.WaitGroup that will be notified when the blaster is done.
func NewBlaster(id string, config *Configuration, wg *sync.WaitGroup) (*Blaster, error) {
	period, err := util.TimeFromFrequency(float64(config.Rate))
	if err != nil {
		return nil, err
	}

	done := make(chan bool)
	return &Blaster{
		id:       id,
		config: config,
		ticker:   time.NewTicker(period),
		httpClient:  &http.Client{},
		wg:       wg,
		done:     done}, nil
}

// Start is a non-blocking call that will start the blaster.
func (b *Blaster) Start() {
	if b.running {
		log.Printf("Blaster %s is already running", b.id)
		return
	}

	go func() {
		time.Sleep(b.config.Duration)
		b.done <- true
	}()

	go run(b)
	log.Printf("Blaster %s started", b.id)
}

// Signal stop to the blaster.
func (b *Blaster) Stop() {
	if !b.running {
		log.Printf("Blaster %s is not running", b.id)
		return
	}

	go func() {
		b.done <- true
	}()
	log.Printf("Blaster %s was signaled to stop", b.id)
}

// TotalRequests returns the current count of the total requests sent.
func (b *Blaster) TotalRequests() int {
	return b.totalRequests
}

// SuccessfulRequests returns the current count of the number of
// successful requests sent. A successful request has a response
// status code less than 400.
func (b *Blaster) SuccessfulRequests() int {
	return b.successfulRequests
}

func run(b *Blaster) {
	b.running = true
	defer b.wg.Done()
	defer b.ticker.Stop()

	for {
		select {
		case <-b.done:
			b.running = false
			return
		case <-b.ticker.C:
			res, err := b.send()
			if err != nil {
				log.Fatalf("Blaster %s failed during send: %T: %v", b.id, err, err)
			}

			if res.StatusCode < 400 {
				b.successfulRequests++
			}

			b.totalRequests++
		}
	}
}

func (b *Blaster) send() (*http.Response, error) {
	req, err := b.config.BuildRequest()
	if err != nil {
		return nil, err
	}

	start := time.Now()
	res, err := b.httpClient.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"%s %s: %s (%v ms)",
		req.Method,
		req.URL.String(),
		res.Status,
		elapsed)
	return res, err
}
