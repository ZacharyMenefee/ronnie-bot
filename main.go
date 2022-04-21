package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

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
	if environment != DEV && environment != PROD {
		panic(environment + " not a valid environment")
	}

	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}

	// send notification in dev channel if rolling out new prod deployment
	// this works because we're using the dumbest possible deployment
	// mechanism with a single machine :)
	if environment == PROD {
		msg, err := startupMessage()
		if err != nil {
			log.Printf("failed generating startup message: %v", err)
		}
		log.Printf("sending startup message: %s", msg)
		session.ChannelMessageSend(devChannelID, msg)
	}
	session.AddHandler(filterEnvironment(environment, indexHandler))
	if err := session.Open(); err != nil {
		panic(err)
	}

	log.Println("running ronnie bot")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("spinning down ronnie bot")
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
			description: "shows an index of available commands",
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

func startupMessage() (string, error) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return "", fmt.Errorf("not in git repo, cannot find hash")
	}

	gitArgs := []string{"log", "--name-status", "HEAD^..HEAD"}
	out, err := exec.Command("git", gitArgs...).Output()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("**ronnie-bot has been reborn!**\n```\n%s\n```", out), nil
}
