package utils

import (
	"os"
	"os/exec"
	"strings"
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
