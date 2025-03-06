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
			if strings.Contains(line, "#msg") || strings.Contains(line, "joined the game") || strings.Contains(line, "left the game") {
	            msg := strings.Split(line, ":")
		        s.ChannelMessageSend(m.ChannelID, msg[len(msg)-1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading log:", err)
	}

	cmd.Wait()
}
