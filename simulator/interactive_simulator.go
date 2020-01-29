package simulator

import (
	"bufio"
	"errors"
	"fmt"
	numeron "github.com/motoki317/numeron-solver"
	"os"
	"regexp"
	"strconv"
)

type InteractiveSimulator struct {
	reader     *bufio.Reader
	re         *regexp.Regexp
	charLength int
}

func NewInteractiveSimulator(charLength int) *InteractiveSimulator {
	reader := bufio.NewReader(os.Stdin)
	reStr := fmt.Sprintf("([0-9]{%v}) ([0-%v]) ([0-%v])", charLength, charLength, charLength)
	// fmt.Println("Debug: regex: ", reStr)
	re := regexp.MustCompile(reStr)

	return &InteractiveSimulator{
		reader:     reader,
		re:         re,
		charLength: charLength,
	}
}

func (s *InteractiveSimulator) AskAnswer(states []*numeron.State) (*numeron.Answer, error) {
	fmt.Println("Type in \"end\" to end and show result.")
	fmt.Print("Next answer in format of \"answer hit blow\": ")
	text, err := s.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if match := s.re.FindStringSubmatch(text); match != nil && len(match) == 4 && len(match[1]) == s.charLength {
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
