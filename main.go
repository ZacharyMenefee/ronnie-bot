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
	bgsID        = "214507567748087808"
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
	if err := session.Open(); err != nil {
		panic(err)
	}
	defer session.Close()

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
	if err := registerActions(session, actions()); err != nil {
		log.Printf("failed to register actions: %v", err)
	}

	log.Println("running ronnie bot")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

type msgHandler func(args []string) string

type action struct {
	name, description string
	handler           msgHandler
}

func actions() []action {
	return []action{
		{
			name:        "help",
			description: "shows an index of available commands",
			handler:     helpHandler,
		},
		{
			name:        "sleepaway",
			description: "gives a random quote from sleepaway camp",
			handler:     sleepawayHandler,
		},
	}
}

// registerActions register our internal representation of text actions as slash commands.
func registerActions(session *discordgo.Session, actions []action) error {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	for _, action := range actions {
		action := action
		cmd := &discordgo.ApplicationCommand{
			Name:        action.name,
			Description: action.description,
		}
		_, err := session.ApplicationCommandCreate(session.State.User.ID, bgsID, cmd)
		if err != nil {
			return fmt.Errorf("failed to create application command %w", err)
		}

		commandHandlers[action.name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// TODO(Monkeyanator) implement command arguments.
			response := action.handler([]string{})
			log.Printf("Handling action %s", action.name)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		}
	}
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	return nil
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
