package client

import numeron "github.com/motoki317/numeron-solver"

// first client always picks the first possible combination when sorted by dictionary.
type FirstClient struct{}

func (c *FirstClient) Query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	return states[0], nil
}
