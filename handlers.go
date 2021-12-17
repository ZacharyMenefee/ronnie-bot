package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	_ "embed"
)

// indexHandler is the entry point that takes all registered actions
// and dispatches to them.
func indexHandler(s *discordgo.Session, mc *discordgo.MessageCreate) {
	// ignore messages from self
	if s.State.User.ID == mc.Author.ID {
		return
	}

	message := mc.Message.Content
	for _, action := range actions() {
		if strings.HasPrefix(message, action.prefix) {
			log.Printf("prefix: %q, matched in message: %s", action.prefix, message)
			args := strings.Split(message, " ")[1:]
			response := action.handler(args)
			s.ChannelMessageSend(mc.ChannelID, prepare(response))
		}
	}
}

// helpHandler provides an index of all registered actions.
func helpHandler(args []string) string {
	sb := strings.Builder{}
	for _, action := range actions() {
		sb.WriteString(fmt.Sprintf("**%s**: %s\n", action.prefix, action.description))
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
