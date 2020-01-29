package client

import (
	numeron "github.com/motoki317/numeron-solver"
	"math/rand"
	"time"
)

// random client randomly picks the next move.
type RandomClient struct {
	rand *rand.Rand
}

func NewRandomClient() *RandomClient {
	return &RandomClient{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *RandomClient) Query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	i := c.rand.Intn(len(states))
	return states[i], nil
}
