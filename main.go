package main

import (
	"flag"
	"log"
	"minedc/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	GuildID        *string
	BotToken       *string
	RemoveCommands *bool
	s              *discordgo.Session
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

func init() {
	GuildID = flag.String("guild", os.Getenv("GUILD_ID"), "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", os.Getenv("BOT_TOKEN"), "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	flag.Parse()

	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	integerOptionMinValue         = 1.0
	dmPermission                  = false
	defaultMemberPermisions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start Minecrfat Server",
		},
		{
			Name:        "stop",
			Description: "Stop Minecrfat Server",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**Starting Minecraft Server...**",
				},
			})
			utils.MessageHandler("start", s, i)
		},
		"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**Stopping Minecraft server...**",
				},
			})
			utils.MessageHandler("stop", s, i)
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		commandsGlo, err := s.ApplicationCommands(s.State.User.ID, "")
		if err != nil {
			log.Fatalf("Cannot fetch guild commands: %v", err)
		}
		commandsGui, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		if err != nil {
			log.Fatalf("Cannot fetch guild commands: %v", err)
		}
		log.Printf("CMD: %v - %v", commandsGlo, commandsGui)
		for _, cmd := range commandsGlo {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
			if err != nil {
				log.Printf("Cannot delete global command %s: %v", cmd.Name, err)
			} else {
				log.Printf("Deleted global command: %s", cmd.Name)
			}
		}
		for _, cmd := range commandsGui {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, cmd.ID)
			if err != nil {
				log.Printf("Cannot delete guild command %s: %v", cmd.Name, err)
			} else {
				log.Printf("Deleted guild command: %s", cmd.Name)
			}
		}

	}

	log.Println("Gracefully shutting down.")
}
