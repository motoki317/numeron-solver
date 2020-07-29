package main

import (
	"fmt"
	numeron "github.com/motoki317/numeron-solver"
	"github.com/motoki317/numeron-solver/client"
	"github.com/motoki317/numeron-solver/simulator"
	"math/rand"
	"strings"
	"time"
)

const (
	charLength      = 4
	charSet         = "0123456789"
	allowDuplicates = true
)

func main() {
	// game()
	supportSolving()
	// simulateGame(1000, false)
}

// Manually plays the game.
func game() {
	c := client.NewInteractiveClient(charLength)
	simulateGameOnce(c, true)
}

// (When we do not know the answer) Interactively records answers and supports solving the game.
func supportSolving() {
	sim := simulator.NewInteractiveSimulator(charLength)
	solver := numeron.NewSolver(charLength, []rune(charSet), allowDuplicates)

	for len(solver.States) > 1 {
		printStates(solver.States)
		ans, err := sim.AskAnswer()
		if err != nil {
			fmt.Println(err)
		} else {
			solver.RecordAnswer(*ans)
		}
	}

	printStates(solver.States)
	fmt.Printf("Finished, took %v total attempts", solver.Attempts)
}

// Simulates the game.
func simulateGames(count int, debug bool) {
	totalAttempts := make(map[string]int)
	clients := make(map[string]client.Client)

	clients["random"] = client.NewRandomClient()
	clients["first"] = &client.FirstClient{}
	clients["different"] = client.NewDifferentClient()

	for k, v := range clients {
		totalAttempts[k] = 0
		for i := 0; i < count; i++ {
			totalAttempts[k] += simulateGameOnce(v, debug)
		}
	}

	for k := range clients {
		fmt.Printf("%v client: %v attempts were made on average.\n", k, float64(totalAttempts[k])/float64(count))
	}
}

// simulate the game once and returns how many attempts it took.
func simulateGameOnce(client client.Client, debug bool) int {
	solver := numeron.NewSolver(charLength, []rune(charSet), allowDuplicates)
	// pick a random state for answer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	answer := solver.States[r.Intn(len(solver.States))]
	if debug {
		printStates(solver.States)
	}

	// simulate the game
	previous := make([]*numeron.Answer, 0)
	for {
		picked, _ := client.Query(solver.States, previous)
		result := answer.Check(picked)
		solver.RecordAnswer(result)
		previous = append(previous, &result)

		if debug {
			fmt.Printf("%v-th call: %v, %v hit, %v blow\n",
				solver.Attempts, string(result.Query.Numbers), result.Hit, result.Blow)
			printStates(solver.States)
		}

		if string(picked.Numbers) == string(answer.Numbers) {
			break
		}
	}

	if debug {
		fmt.Printf("%v total attempts were made.\n", solver.Attempts)
	}
	return solver.Attempts
}

func printStates(states []*numeron.State) {
	arr := make([]string, 0)
	if len(states) < 100 {
		for _, state := range states {
			arr = append(arr, string(state.Numbers))
		}
	}
	fmt.Println(strings.Join(arr, ", "))
	fmt.Printf("%v possible states remaining\n", len(states))
}
