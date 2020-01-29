package client

import numeron "github.com/motoki317/numeron-solver"

// Client plays the game
type Client interface {
	Query(current []*numeron.State, previous []*numeron.Answer) (*numeron.State, error)
}
