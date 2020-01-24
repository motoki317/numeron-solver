package main

import (
	"bufio"
	"errors"
	"fmt"
	numeron "github.com/motoki317/numeron-solver"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	charLength = 3
	charSet    = "0123456789"
)

func main() {
	simulateGame(1000, false)
}

// Manually plays the game.
func game() {
	c := newInteractiveClient()
	simulateGameOnce(c, true)
}

// (When we do not know the answer) Interactively records answers and supports solving the game.
func supportSolving() {
	sim := newInteractiveSimulator()
	solver := numeron.NewSolver(charLength, []rune(charSet))

	for len(solver.States) > 1 {
		ans, err := sim.askAnswer(solver.States)
		if err != nil {
			fmt.Println(err)
		} else {
			solver.RecordAnswer(*ans)
		}
	}

	fmt.Printf("Finished, took %v total attempts", solver.Attempts)
}

// Simulates the game.
func simulateGame(count int, debug bool) {
	totalAttempts := make(map[string]int)
	clients := make(map[string]client)

	clients["random"] = newRandomClient()
	clients["first"] = &firstClient{}
	clients["different"] = newDifferentClient()

	for k, v := range clients {
		totalAttempts[k] = 0
		for i := 0; i < count; i++ {
			totalAttempts[k] += simulateGameOnce(v, debug)
		}
	}

	for k := range clients {
		fmt.Printf("%v client: %v attempts were made.\n", k, float64(totalAttempts[k])/float64(count))
	}
}

// simulate the game once and returns how many attempts it took.
func simulateGameOnce(client client, debug bool) int {
	solver := numeron.NewSolver(charLength, []rune(charSet))
	// pick a random state for answer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	answer := solver.States[r.Intn(len(solver.States))]
	if debug {
		printStates(solver.States)
	}

	// simulate the game
	previous := make([]*numeron.Answer, 0)
	for {
		picked, _ := client.query(solver.States, previous)
		result := answer.Check(picked)
		solver.RecordAnswer(result)
		previous = append(previous, &result)

		if debug {
			fmt.Printf("%v-th call: %v, %v hit, %v blow",
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

// When we do not know the answer and want to know all possible remaining combinations
type simulator interface {
	askAnswer([]*numeron.State) (*numeron.Answer, error)
}

type interactiveSimulator struct {
	reader *bufio.Reader
	re     *regexp.Regexp
}

// Client plays the game
type client interface {
	query(current []*numeron.State, previous []*numeron.Answer) (*numeron.State, error)
}

type interactiveClient struct {
	reader *bufio.Reader
	re     *regexp.Regexp
}

// random client randomly picks the next move.
type randomClient struct {
	rand *rand.Rand
}

// first client always picks the first possible combination when sorted by dictionary.
type firstClient struct{}

// this client always tries to pick unused numbers if possible.
type differentClient struct {
	rand *rand.Rand
}

func newInteractiveSimulator() *interactiveSimulator {
	reader := bufio.NewReader(os.Stdin)
	reStr := fmt.Sprintf("([0-9]{%v}) ([0-%v]) ([0-%v])", charLength, charLength, charLength)
	fmt.Println("Debug: regex: ", reStr)
	re := regexp.MustCompile(reStr)

	return &interactiveSimulator{
		reader: reader,
		re:     re,
	}
}

func (s *interactiveSimulator) askAnswer(states []*numeron.State) (*numeron.Answer, error) {
	printStates(states)

	fmt.Println("Type in \"end\" to end and show result.")
	fmt.Print("Next answer in format of \"answer hit blow\": ")
	text, err := s.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if match := s.re.FindStringSubmatch(text); match != nil && len(match) == 4 && len(match[1]) == charLength {
		hit, err := strconv.Atoi(match[2])
		if err != nil {
			fmt.Println("Invalid format for hit number.")
		}
		blow, err := strconv.Atoi(match[3])
		if err != nil {
			fmt.Println("Invalid format for blow number.")
		}

		return &numeron.Answer{
			Query: numeron.State{
				Numbers: []rune(match[1]),
			},
			Hit:  hit,
			Blow: blow,
		}, nil
	} else {
		fmt.Printf("Debug: matched: %v\n", match)
		return nil, errors.New("invalid input")
	}
}

func newInteractiveClient() *interactiveClient {
	reader := bufio.NewReader(os.Stdin)
	reStr := fmt.Sprintf("([0-9]{%v})", charLength)
	fmt.Println("Debug: regex: ", reStr)
	re := regexp.MustCompile(reStr)

	return &interactiveClient{
		reader: reader,
		re:     re,
	}
}

func (c *interactiveClient) query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	fmt.Print("Next answer: ")
	text, err := c.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if match := c.re.FindStringSubmatch(text); match != nil && len(match) == 2 && len(match[1]) == charLength {
		return &numeron.State{
			Numbers: []rune(match[1]),
		}, nil
	} else {
		fmt.Printf("Debug: matched: %v\n", match)
		return nil, errors.New("invalid input")
	}
}

func newRandomClient() *randomClient {
	return &randomClient{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func newDifferentClient() *differentClient {
	return &differentClient{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *randomClient) query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	i := c.rand.Intn(len(states))
	return states[i], nil
}

func (c *firstClient) query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	return states[0], nil
}

func (c *differentClient) query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
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

func printStates(states []*numeron.State) {
	fmt.Println()
	if len(states) < 100 {
		for _, state := range states {
			fmt.Println(string(state.Numbers))
		}
	}
	fmt.Printf("%v possible states remaining\n", len(states))
}
