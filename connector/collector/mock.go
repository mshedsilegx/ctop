//go:build !release

package collector

import (
	"math/rand"
	"time"

	"github.com/bcicen/ctop/models"
)

// Mock collector
type Mock struct {
	models.Metrics
	stream     chan models.Metrics
	done       bool
	running    bool
	aggression int64
}

func NewMock(a int64) *Mock {
	c := &Mock{
		Metrics:    models.Metrics{},
		aggression: a,
	}
	c.MemLimit = 2147483648
	return c
}

func (c *Mock) Running() bool {
	return c.running
}

func (c *Mock) Start() {
	c.done = false
	c.stream = make(chan models.Metrics)
	go c.run()
}

func (c *Mock) Stop() {
	c.running = false
	c.done = true
}

func (c *Mock) Stream() chan models.Metrics {
	return c.stream
}

func (c *Mock) Logs() LogCollector {
	return &MockLogs{make(chan bool)}
}

func (c *Mock) run() {
	c.running = true
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	defer close(c.stream)

	// set to random static value, once
	c.Pids = r.Intn(12)
	c.IOBytesRead = r.Int63n(8098) * c.aggression
	c.IOBytesWrite = r.Int63n(8098) * c.aggression

	for {
		c.CPUUtil += r.Intn(2) * int(c.aggression)
		if c.CPUUtil >= 100 {
			c.CPUUtil = 0
		}

		c.NetTx += r.Int63n(60) * c.aggression
		c.NetRx += r.Int63n(60) * c.aggression
		c.MemUsage += r.Int63n(c.MemLimit/512) * c.aggression
		if c.MemUsage > c.MemLimit {
			c.MemUsage = 0
		}
		c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
		c.stream <- c.Metrics
		if c.done {
			break
		}
		time.Sleep(1 * time.Second)
	}

	c.running = false
}
