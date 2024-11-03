package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	// Fetch config from GitHub repository
	config, err := fetchConfigFromGitHub("github_user", "repo", "config.json", "auth")
	if err != nil {
		log.Println("Error fetching config file:", err)
		return
	}

	instructions := Instructions{
		DataURL:  "https://pjreddie.com/media/files/mnist_train.csv",
		ModelURL: "",
		Command:  "python src/train.py --data_url https://pjreddie.com/media/files/mnist_train.csv",
	}

	// Send the instructions to the public URL
	if err := sendInstructions(config.PublicURL, instructions); err != nil {
		log.Println("Error sending instructions:", err)
		return
	}

	log.Println("Instructions sent successfully!")
}

func fetchConfigFromGitHub(owner, repo, filePath, token string) (Config, error) {
	var config Config

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, filePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return config, fmt.Errorf("error creating request: %w", err)
	}

	// Set the authorization header with your GitHub token
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return config, fmt.Errorf("error fetching config from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return config, fmt.Errorf("failed to fetch config: %d, Response: %s", resp.StatusCode, body)
	}

	var githubResponse struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubResponse); err != nil {
		return config, fmt.Errorf("error decoding GitHub response: %w", err)
	}

	// Decode the base64 encoded content
	configData, err := base64.StdEncoding.DecodeString(githubResponse.Content)
	if err != nil {
		return config, fmt.Errorf("error decoding base64 content: %w", err)
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		return config, fmt.Errorf("error unmarshalling config JSON: %w", err)
	}

	return config, nil
}

func sendInstructions(publicURL string, instructions Instructions) error {
	jsonData, err := json.Marshal(instructions)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/execute", publicURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending instructions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to send instructions: %d, Response: %s", resp.StatusCode, body)
	}

	log.Println("Instructions sent successfully!")
	return nil
}
