package helpers

import (
	"fmt"
	"io"
	"net/http"
)

func Contains(s string, slice []string) bool {
	found := false
	for _, e := range slice {
		if e == s {
			found = true
			break
		}
	}

	return found
}

func GetIp() (string, error) {
	ipServiceURL := "https://ifconfig.me"
	resp, err := http.Get(ipServiceURL)
	if err != nil {
		return "", fmt.Errorf("Error making HTTP request: %v\n", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to get IP: %v\n", err)
	}

	return string(body), nil
}
