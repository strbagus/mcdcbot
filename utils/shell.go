package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Systemctl(arg string) {
	servName := os.Getenv("SERVICE_NAME")
	cmd := exec.Command("sudo", "systemctl", arg, servName)
	cmd.Run()
}
func IsServiceRunning(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()

	if err != nil {
		return false
	}

	status := strings.TrimSpace(string(output))
	return status == "active"
}

func LogListen(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	servName := os.Getenv("SERVICE_NAME")
	cmd := exec.CommandContext(ctx, "journalctl", "-u", servName, "-f")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping log listening...")
			return
		default:
			line := scanner.Text()
			fmt.Println("server: ", line)
			if strings.Contains(line, "#msg") || strings.Contains(line, "joined the game") || strings.Contains(line, "left the game") || strings.Contains(line, "Done preparing level") {
				msg := strings.Split(line, ":")
				tmp := msg[len(msg)-1]
				var res string
				if strings.Contains(line, "joined the game") {
					name := strings.TrimSuffix(tmp, " joined the game")
					res = fmt.Sprintf("%s: **joined the game**", name)
				} else if strings.Contains(line, "left the game") {
					name := strings.TrimSuffix(tmp, " left the game")
					res = fmt.Sprintf("%s: **left the game**", name)
				} else if strings.Contains(line, "#msg") {
					start := strings.Index(tmp, "<") + 1
					end := strings.Index(tmp, ">")
					username := tmp[start:end]
					message := strings.TrimPrefix(tmp[end+2:], "#msg ")
					res = fmt.Sprintf("%s: *%s*", username, message)
				} else {
                    res = "**Minecraft Server is Running!**"
                }
				s.ChannelMessageSend(m.ChannelID, res)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading log:", err)
	}

	cmd.Wait()
}
