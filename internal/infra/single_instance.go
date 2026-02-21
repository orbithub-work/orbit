package infra

import (
	"fmt"
	"net/http"
	"time"
)

const CorePort = 32000

// CheckAndNotifyExistingInstance checks if another instance is running.
// If running, it notifies it to show the main window and returns true.
func CheckAndNotifyExistingInstance() bool {
	// 1. Try to connect to the port
	url := fmt.Sprintf("http://127.0.0.1:%d/api/internal/show_main", CorePort)
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	// Use POST as required by the handler
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		// Connection failed, likely no instance running
		return false
	}
	defer resp.Body.Close()

	// If we get a response, it means something is listening.
	// We assume it is our application.
	return true
}
