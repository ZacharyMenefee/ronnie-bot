package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	_ "embed"
)

// helpHandler provides an index of all registered actions.
func helpHandler(args []string) string {
	sb := strings.Builder{}
	for _, action := range actions() {
		sb.WriteString(fmt.Sprintf("**%s**: %s\n", action.name, action.description))
	}
	return sb.String()
}

//go:embed static/camp.txt
var script string

// sleepawayHandler prints a random quote from the Sleepaway Camp movie
// TODO(Monkeyanator) extend to sending a screencap.
func sleepawayHandler(args []string) string {
	quotes := strings.Split(script, "\n")
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("> %s", quotes[r.Int()&len(quotes)])
}

func prepare(input string) string {
	return strings.TrimSuffix(input, "\n")
}
