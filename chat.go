package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"llama-chat-cli/src"
	"net/http"
	"time"
)

type Post struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string      `json:"role"`
			Content string      `json:"content"`
			Refusal interface{} `json:"refusal"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens            int `json:"prompt_tokens"`
		CompletionTokens        int `json:"completion_tokens"`
		TotalTokens             int `json:"total_tokens"`
		CompletionTokensDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func main() {
	configPath := "./config.yaml"
	config, err := ingestyaml.Ingest(configPath)
	if err != nil {
		fmt.Println(err)
	}
	endpoint := config.Api.ChatEndpoint
	key := config.Api.Key
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}
	umessage := "Hello how are you today?"
	chat := "gpt-4o-mini"

	pReq := Post{
		Model: chat,
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: umessage},
		},
	}

	body, err := json.Marshal(pReq)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("Header-Name", "Header-Value")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	Rbody, err := io.ReadAll(resp.Body)
	fmt.Println(string(Rbody))
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	response := ChatResponse{}
	err = json.Unmarshal(Rbody, &response)
	if err != nil {
		fmt.Println("ERRRORRR")
	}
	choices := response.Choices
	formatted := fmt.Sprintf(`
      %s >> "%s"
      Debug: total tokens: %d
      `, choices[0].Message.Role, choices[0].Message.Content, response.Usage.TotalTokens)
	fmt.Println(formatted)
}
