package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	// "regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Systemctl(arg string, s *discordgo.Session, chName string) {
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

func LogListen(ctx context.Context, s *discordgo.Session, channelID string) {
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
	log.Println("Listening Minecraft Server Log..")
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping log listening...")
			return
		default:
			line := scanner.Text()
			// fmt.Println("server: ", line)
			if strings.Contains(line, "#msg") || strings.Contains(line, "joined the game") || strings.Contains(line, "left the game") || strings.Contains(line, "Done preparing level") || strings.Contains(line, "UUID of player") || strings.Contains(line, "You are not whitelisted on this server!") {
			fmt.Println("server: ", line)
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
				/* } else if strings.Contains(line, "UUID of player") {
					re := regexp.MustCompile(`UUID of player (.+?) is ([0-9a-fA-F-]+)`)
					matches := re.FindStringSubmatch(line)
					if len(matches) == 3 {
						name := matches[1]
						uuid := matches[2]
            res = fmt.Sprintf("`{\"uuid\": \"%s\", \"name\": \"%s\"}`", uuid, name)
					}
				} else if strings.Contains(line, "You are not whitelisted on this server!") {
          res = "*User tersebut belum ada di whitelist*" */
				} else if strings.Contains(line, "Done preparing level") {
					res = "**Minecraft Server is Running!**"
					s.ChannelEdit(channelID, &discordgo.ChannelEdit{
						Name: "minecraft-on",
					})
				}
				s.ChannelMessageSend(channelID, res)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading log:", err)
	}

	cmd.Wait()
}
