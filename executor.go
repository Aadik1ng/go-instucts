package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Instructions struct {
	DataURL  string `json:"data_url"`
	ModelURL string `json:"model_url"`
	Command  string `json:"command"`
}

type Config struct {
	IP        string `json:"ip"`
	Port      string `json:"port"`
	Timestamp string `json:"timestamp"`
	PublicURL string `json:"public_url"`
}

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func updateConfig(publicURL string) {
	ip, err := getOutboundIP()
	if err != nil {
		log.Println("Error getting IP address:", err)
		return
	}

	config := Config{
		IP:        ip,
		Port:      "8080",
		Timestamp: time.Now().Format(time.RFC3339),
		PublicURL: publicURL,
	}

	configData, err := json.Marshal(config)
	if err != nil {
		log.Println("Error marshalling config JSON:", err)
		return
	}

	err = os.WriteFile("config.json", configData, 0644)
	if err != nil {
		log.Println("Error writing config file:", err)
		return
	}

	repoPath := "path/to/repo"
	err = os.Chdir(repoPath)
	if err != nil {
		log.Println("Error changing directory:", err)
		return
	}

	if err := pullLatestChanges(); err != nil {
		log.Println("Error pulling latest changes:", err)
		return
	}

	commands := [][]string{
		{"git", "add", "config.json"},
		{"git", "commit", "-m", "Update config with public URL and timestamp"},
		{"git", "push"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Error running command %v: %v\n", args, err)
			return
		}
	}

	log.Println("Config update and push completed successfully")
}

func pullLatestChanges() error {
	cmd := exec.Command("git", "pull", "origin", "main")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error pulling latest changes: %v. Output: %s\n", err, output)
		return err
	}

	log.Println("Pulled latest changes successfully")
	return nil
}

func startNgrok() (string, error) {
	cmd := exec.Command("ngrok", "http", "8080")

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to start ngrok: %w", err)
	}

	time.Sleep(2 * time.Second)

	publicURL, err := getNgrokPublicURL()
	if err != nil {
		return "", err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("ngrok exited with error: %v\nOutput: %s\n", err, output.String())
		} else {
			log.Printf("ngrok output: %s\n", output.String())
		}
	}()

	log.Println("ngrok started successfully with public URL:", publicURL)
	return publicURL, nil
}

func getNgrokPublicURL() (string, error) {
	cmd := exec.Command("curl", "-s", "http://localhost:4040/api/tunnels")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ngrok public URL: %w", err)
	}

	type Tunnel struct {
		PublicURL string `json:"public_url"`
	}
	type NgrokResponse struct {
		Tunnels []Tunnel `json:"tunnels"`
	}

	var ngrokResp NgrokResponse
	if err := json.Unmarshal(output, &ngrokResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal ngrok response: %w", err)
	}

	if len(ngrokResp.Tunnels) == 0 {
		return "", fmt.Errorf("no tunnels found")
	}

	return ngrokResp.Tunnels[0].PublicURL, nil
}

func executePipeline(w http.ResponseWriter, r *http.Request) {
	var instructions Instructions

	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received instructions: %+v\n", instructions)

	go func() {
		cmd := exec.Command("cmd", "/C", instructions.Command)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output

		if err := cmd.Run(); err != nil {
			log.Printf("Error executing command: %v. Output: %s\n", err, output.String())
			return
		}

		log.Printf("Command executed successfully: %s\n", output.String())
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func main() {
	publicURL, err := startNgrok()
	if err != nil {
		log.Println("Error starting ngrok:", err)
		return
	}

	updateConfig(publicURL)

	http.HandleFunc("/execute", executePipeline)
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
