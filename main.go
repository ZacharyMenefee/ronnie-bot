package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	_ "embed"
)

const (
	DEV          = "DEV"
	PROD         = "PROD"
	devChannelID = "845749844437893121"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic("did not receive bot token")
	}
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" || (environment != DEV && environment != PROD) {
		panic("must set valid environment")
	}

	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}

	session.AddHandler(filterEnvironment(environment, indexHandler))
	if err := session.Open(); err != nil {
		panic(err)
	}

	log.Println("running ronnie bot")
	for {
	}
}

type msgHandler func(args []string) string

type action struct {
	prefix, description string
	handler             msgHandler
}

func actions() []action {
	return []action{
		{
			prefix:      "!help",
			description: "provides an index of available commands",
			handler:     helpHandler,
		},
		{
			prefix:      "!sleepaway",
			description: "gives a random quote from sleepaway camp",
			handler:     sleepawayHandler,
		},
	}
}

// filterEnvironment returns a new handler that ignores messages based on the instance's environment.
func filterEnvironment(environment string, handler func(s *discordgo.Session, mc *discordgo.MessageCreate)) func(s *discordgo.Session, mc *discordgo.MessageCreate) {
	return func(s *discordgo.Session, mc *discordgo.MessageCreate) {
		if environment == DEV {
			if mc.ChannelID != devChannelID {
				return
			}
		} else {
			if mc.ChannelID == devChannelID {
				return
			}
		}
		handler(s, mc)
	}
}

func indexHandler(s *discordgo.Session, mc *discordgo.MessageCreate) {
	// ignore messages from self
	if s.State.User.ID == mc.Author.ID {
		return
	}

	message := mc.Message.Content
	log.Printf("message: %s", message)
	for _, action := range actions() {
		if strings.HasPrefix(message, action.prefix) {
			log.Printf("prefix: %q, matched in message: %s", action.prefix, message)
			args := strings.Split(message, " ")[1:]
			response := action.handler(args)
			s.ChannelMessageSend(mc.ChannelID, prepare(response))
		}
	}
}

func helpHandler(args []string) string {
	sb := strings.Builder{}
	for _, action := range actions() {
		sb.WriteString(fmt.Sprintf("**%s**: %s\n", action.prefix, action.description))
	}
	return sb.String()
}

//go:embed static/camp.txt
var script string

func sleepawayHandler(args []string) string {
	quotes := strings.Split(script, "\n")
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("> %s", quotes[r.Int()&len(quotes)])
}

func prepare(input string) string {
	return strings.TrimSuffix(input, "\n")
}
