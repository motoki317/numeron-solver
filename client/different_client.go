package client

import (
	numeron "github.com/motoki317/numeron-solver"
	"math/rand"
	"time"
)

// Different client prefers to pick unused numbers if possible.
type DifferentClient struct {
	rand *rand.Rand
}

func NewDifferentClient() *DifferentClient {
	return &DifferentClient{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *DifferentClient) Query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	used := make(map[rune]bool)
	for _, p := range previous {
		for _, r := range p.Query.Numbers {
			used[r] = true
		}
	}

	calcRank := func(state *numeron.State) (rank int) {
		for _, r := range state.Numbers {
			if isUsed, ok := used[r]; !isUsed || !ok {
				rank++
			}
		}
		return
	}

	maxRank := -1
	var picked numeron.State
	for _, state := range states {
		rank := calcRank(state)
		if maxRank < rank {
			maxRank = rank
			picked = *state
		}
	}

	return &picked, nil
}
