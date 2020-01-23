package main

import (
	"bufio"
	"fmt"
	numeron "github.com/motoki317/numeron-solver"
	"os"
	"regexp"
	"strconv"
)

const (
	charLength = 4
	charSet = "0123456789"
)

func main() {
	solve()
}

func solve() {
	solver := numeron.NewSolver(charLength, []rune(charSet))
	reader := bufio.NewReader(os.Stdin)
	reStr := fmt.Sprintf("([0-9]{%v}) ([0-%v]) ([0-%v])", charLength, charLength, charLength)
	fmt.Println("Debug: regex: ", reStr)
	re := regexp.MustCompile(reStr)

	for len(solver.States) > 1 {
		printSolver(&solver)

		fmt.Println("Type in \"end\" to end and show result.")
		fmt.Print("Next answer in format of \"answer hit blow\": ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if match := re.FindStringSubmatch(text); match != nil && len(match) == 4 && len(match[1]) == charLength {
			hit, err := strconv.Atoi(match[2])
			if err != nil {
				fmt.Println("Invalid format for hit number.")
			}
			blow, err := strconv.Atoi(match[3])
			if err != nil {
				fmt.Println("Invalid format for blow number.")
			}

			ans := numeron.Answer{
				Query: numeron.State{
					Numbers: []rune(match[1]),
				},
				Hit:   hit,
				Blow:  blow,
			}
			solver.RecordAnswer(ans)
		} else {
			fmt.Printf("Debug: matched: %v\n", match)
			fmt.Println("Invalid format.")
		}
	}

	printSolver(&solver)
	fmt.Printf("Ended. %v total attempts was made.\n", solver.Attempts)
}

func validate() {
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

func printSolver(s *numeron.Solver) {
	fmt.Println()
	if len(s.States) < 100 {
		for _, state := range s.States {
			fmt.Println(string(state.Numbers))
		}
	}
	fmt.Printf("%v states remaining\n", len(s.States))
}
