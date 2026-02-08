package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const blenderAddonPort = 29877

// BlenderSendRequest is the payload sent to the Blender addon.
type BlenderSendRequest struct {
	Action string   `json:"action"` // "import"
	Files  []string `json:"files"`  // absolute paths
}

// BlenderStatus represents the connection status with the Blender addon.
type BlenderStatus struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

// SendToBlender sends one or more file paths to the Blender addon for import.
func (a *App) SendToBlender(absolutePaths []string) BlenderStatus {
	payload := BlenderSendRequest{
		Action: "import",
		Files:  absolutePaths,
	}
	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("http://127.0.0.1:%d/sushi", blenderAddonPort)

	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return BlenderStatus{Connected: false, Error: "Blender addon not running. Open Blender and enable the Sushi addon."}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return BlenderStatus{Connected: true, Error: fmt.Sprintf("Blender returned %d: %s", resp.StatusCode, string(respBody))}
	}

	return BlenderStatus{Connected: true}
}

// PingBlender checks if the Blender addon is reachable.
func (a *App) PingBlender() BlenderStatus {
	client := &http.Client{Timeout: 1 * time.Second}
	url := fmt.Sprintf("http://127.0.0.1:%d/sushi", blenderAddonPort)

	resp, err := client.Get(url)
	if err != nil {
		return BlenderStatus{Connected: false, Error: "Blender addon not reachable"}
	}
	defer resp.Body.Close()

	return BlenderStatus{Connected: true}
}
