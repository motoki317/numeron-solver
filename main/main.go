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
	count := 10000
	totalAttemptsRandom := 0
	totalAttemptsFirst := 0
	debug := false

	rc := newRandomClient()
	fc := &firstClient{}

	for i := 0; i < count; i++ {
		totalAttemptsRandom += simulateGameOnce(rc, debug)
	}
	for i := 0; i < count; i++ {
		totalAttemptsFirst += simulateGameOnce(fc, debug)
	}

	fmt.Printf("Random client: %v attempts were made.\n", float64(totalAttemptsRandom)/float64(count))
	fmt.Printf("First client: %v attempts were made.\n", float64(totalAttemptsFirst)/float64(count))
}

// simulate the game once and returns how many attempts it took.
func simulateGameOnce(client client, debug bool) int {
	solver := numeron.NewSolver(charLength, []rune(charSet))
	// pick a random state for answer
	rc := newRandomClient()
	answer, _ := rc.query(solver.States)
	if debug {
		printStates(solver.States)
	}

	// simulate the game
	for {
		picked, _ := client.query(solver.States)
		result := answer.Check(picked)
		solver.RecordAnswer(result)

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
	recordAnswer([]*numeron.State) (*numeron.Answer, error)
}

type interactiveSimulator struct {
	reader *bufio.Reader
	re     *regexp.Regexp
}

// Client plays the game
type client interface {
	query([]*numeron.State) (*numeron.State, error)
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

func (s *interactiveSimulator) recordAnswer(states []*numeron.State) (*numeron.Answer, error) {
	printStates(states)

	fmt.Println("Type in \"end\" to end and show result.")
	fmt.Print("Next recordAnswer in format of \"answer hit blow\": ")
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

func (c *interactiveClient) query(states []*numeron.State) (*numeron.State, error) {
	printStates(states)

	fmt.Println("Type in \"end\" to end and show result.")
	fmt.Print("Next recordAnswer in format of \"answer hit blow\": ")
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

func (c *randomClient) query(states []*numeron.State) (*numeron.State, error) {
	i := c.rand.Intn(len(states))
	return states[i], nil
}

func (c *firstClient) query(states []*numeron.State) (*numeron.State, error) {
	return states[0], nil
}

func validateInteractive() {
	reader := bufio.NewReader(os.Stdin)
	answer := numeron.State{Numbers: []rune("5460")}
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		res := answer.Check(&numeron.State{
			Numbers: []rune(text),
		})
		fmt.Printf("Hit: %v, Blow: %v", res.Hit, res.Blow)
	}
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
