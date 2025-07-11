package utils

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"os"

	"github.com/bwmarrin/discordgo"
)

var cancelFunc context.CancelFunc

func MessageHandler(command string, s *discordgo.Session, m *discordgo.InteractionCreate) {
	log.Printf("Handle %v.", command)

	servName := os.Getenv("SERVICE_NAME")

	switch command {
	case "start":
		if IsServiceRunning(servName) {
			msg := "Minecraft Server is already Running!"
			s.ChannelMessageSend(m.ChannelID, msg)
			log.Println(msg)
		} else {
			startServer()
			if cancelFunc != nil {
				fmt.Println("Log already running")
				return
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancelFunc = cancel
			go LogListen(ctx, s, m.ChannelID)
		}
	case "stop":
		var msg string
		if IsServiceRunning(servName) {
			msg = "**Minecraft Server is Stopped!**"
			cmd := exec.Command("sudo", "systemctl", "stop", servName)
			cmd.Run()
			log.Println("Server Service is stopped.")
			if cancelFunc != nil {
				cancelFunc()
				cancelFunc = nil
				fmt.Println("Log stopped")
			} else {
				fmt.Println("Log is not running")
			}
		} else {
			msg = "Minecraft Server is not Runnning!"
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		log.Println("Message: ", msg)
		s.ChannelEdit(m.ChannelID, &discordgo.ChannelEdit{
			Name: "minecraft-off",
		})
		log.Println("Channel name updated.")
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
	}
}

func startServer() {
	servName := os.Getenv("SERVICE_NAME")
	cmd := exec.Command("sudo", "systemctl", "start", servName)
	cmd.Run()
	log.Println("Starting Server Service.")
}
