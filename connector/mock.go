//go:build !release

package connector

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/container"
	"github.com/jgautheron/codename-generator"
	"github.com/nu7hatch/gouuid"
)

func init() { enabled["mock"] = NewMock }

type Mock struct {
	containers container.Containers
}

func NewMock() (Connector, error) {
	cs := &Mock{}
	go cs.Init()
	go cs.Loop()
	return cs, nil
}

// Create Mock containers
func (cs *Mock) Init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 4; i++ {
		cs.makeContainer(r, 3, true)
	}

	for i := 0; i < 16; i++ {
		cs.makeContainer(r, 1, false)
	}

}

func (cs *Mock) Wait() struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(30 * time.Second)
		close(ch)
	}()
	return <-ch
}

var healthStates = []string{"starting", "healthy", "unhealthy"}

func (cs *Mock) makeContainer(r *rand.Rand, aggression int64, health bool) {
	collector := collector.NewMock(aggression)
	manager := manager.NewMock()
	c := container.New(makeID(), collector, manager)
	c.SetMeta("name", makeName())
	c.SetState(makeState(r))
	if health {
		var i int
		c.SetMeta("health", healthStates[i])
		go func() {
			for {
				i++
				if i >= len(healthStates) {
					i = 0
				}
				c.SetMeta("health", healthStates[i])
				time.Sleep(12 * time.Second)
			}
		}()
	}
	cs.containers = append(cs.containers, c)
}

func (cs *Mock) Loop() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	iter := 0
	for {
		// Change state for random container
		if iter%5 == 0 && len(cs.containers) > 0 {
			randC := cs.containers[r.Intn(len(cs.containers))]
			randC.SetState(makeState(r))
		}
		iter++
		time.Sleep(3 * time.Second)
	}
}

// Get a single container, by ID
func (cs *Mock) Get(id string) (*container.Container, bool) {
	for _, c := range cs.containers {
		if c.Id == id {
			return c, true
		}
	}
	return nil, false
}

// All returns array of all containers, sorted by field
func (cs *Mock) All() container.Containers {
	cs.containers.Sort()
	cs.containers.Filter()
	return cs.containers
}

func makeID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(u.String(), "-", "")[:12]
}

func makeName() string {
	n, err := codename.Get(codename.Sanitized)
	nsp := strings.Split(n, "-")
	if len(nsp) > 2 {
		n = strings.Join(nsp[:2], "-")
	}
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(n, "-", "_")
}

func makeState(r *rand.Rand) string {
	switch r.Intn(10) {
	case 0, 1, 2:
		return "exited"
	case 3:
		return "paused"
	}
	return "running"
}
