package client

import (
	"bufio"
	"errors"
	"fmt"
	numeron "github.com/motoki317/numeron-solver"
	"os"
	"regexp"
)

type InteractiveClient struct {
	reader     *bufio.Reader
	re         *regexp.Regexp
	charLength int
}

func NewInteractiveClient(charLength int) *InteractiveClient {
	reader := bufio.NewReader(os.Stdin)
	reStr := fmt.Sprintf("([0-9]{%v})", charLength)
	// fmt.Println("Debug: regex: ", reStr)
	re := regexp.MustCompile(reStr)

	return &InteractiveClient{
		reader:     reader,
		re:         re,
		charLength: charLength,
	}
}

func (c *InteractiveClient) Query(states []*numeron.State, previous []*numeron.Answer) (*numeron.State, error) {
	fmt.Print("Next answer: ")
	text, err := c.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if match := c.re.FindStringSubmatch(text); match != nil && len(match) == 2 && len(match[1]) == c.charLength {
		return &numeron.State{
			Numbers: []rune(match[1]),
		}, nil
	} else {
		fmt.Printf("Debug: matched: %v\n", match)
		return nil, errors.New("invalid input")
	}
}
