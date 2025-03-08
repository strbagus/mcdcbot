package utils

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

var cancelFunc context.CancelFunc

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    chId := os.Getenv("CHANNEL_ID")
	if m.Author.ID == s.State.User.ID || chId != m.ChannelID {
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
			msg = "Minecraft Server is already Running!"
		} else {
			msg = "**Starting Minecraft Server...**"
			Systemctl("start")

			if cancelFunc != nil {
				fmt.Println("Log already running")
				return
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancelFunc = cancel
			go LogListen(ctx, s, m)
			s.ChannelEdit(m.ChannelID, &discordgo.ChannelEdit{
				Name: "minecraft-on",
			})
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	case "stop":
		var msg string
		if IsServiceRunning(servName) {
			msg = "**Minecraft Server is Stopped!**"
			Systemctl("stop")
			if cancelFunc != nil {
				cancelFunc()
				cancelFunc = nil
				fmt.Println("Log stopped")
			} else {
				fmt.Println("Log is not running")
			}
			s.ChannelEdit(m.ChannelID, &discordgo.ChannelEdit{
				Name: "minecraft-off",
			})
		} else {
			msg = "Minecraft Server is not Runnning!"
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
	}
}
