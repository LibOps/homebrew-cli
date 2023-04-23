package gcloud

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func AccessToken() (string, error) {
	args := []string{
		"auth",
		"print-identity-token",
	}
	cmd := exec.Command("gcloud", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	log.Println("Generated gcloud identity token")
	token := stdout.String()
	return strings.TrimSpace(token), nil
}
