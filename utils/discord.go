package utils

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	prefix := "!"
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	args := strings.Fields(m.Content[len(prefix):])
	if len(args) == 0 {
		return
	}

	command := args[0]
	servName := os.Getenv("SERVICE_NAME")

	switch command {
	case "start":
		var msg string
		if IsServiceRunning(servName) {
			msg = "Minecraft Server is Active!"
		} else {
			msg = "Start Minecraft Server!"
			Systemctl("start")
			s.ChannelEdit(m.ChannelID, &discordgo.ChannelEdit{
				Name: "minecraft-on",
			})
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	case "stop":
		var msg string
		if IsServiceRunning(servName) {
			msg = "Stop Minecraft Server!"
			Systemctl("stop")
			s.ChannelEdit(m.ChannelID, &discordgo.ChannelEdit{
				Name: "minecraft-off",
			})
		} else {
			msg = "Minecraft Server is Inactive!"
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
	}
}
