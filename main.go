package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Request body to llama API
type chatRequest struct {
	Model        string       `json:"model"`
	Messages     []reqMessage `json:"messages"`
	Functions    []function   `json:"functions"`
	Stream       bool         `json:"stream"`
	FunctionCall string       `json:"function_call"`
}

type reqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  parameters `json:"parameters"`
	Required    []string   `json:"required"`
}

type parameters struct {
	Type       string     `json:"type"`
	Properties properties `json:"properties"`
}

type properties struct {
	Words words `json:"words"`
}

type words struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Response body from llama API
type chatResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Index        int        `json:"index"`
	Message      resMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type resMessage struct {
	Role         string       `json:"role"`
	Content      string       `json:"content"`
	FunctionCall functionCall `json:"function_call"`
}

type functionCall struct {
	Name      string    `json:"name"`
	Arguments arguments `json:"arguments"`
}

type arguments struct {
	Words []string `json:"words"`
}

// Endpoint
const API_URL = "https://api.llama-api.com/chat/completions"

// Get generated response from Llama API
func getGeneratedResponse(prompt string) string {
	//Create request body
	chatReq := createChatRequest(prompt)

	//Marshal Go struct into Json
	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}

	//Create Http request struct with request method, endpoint and request body
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", API_URL, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Failed to create http request struct: %v", err)
	}

	//Add necessary headers, including the API key for authorization
	apiKey := os.Getenv("LLAMA_API_KEY")
	if apiKey == "" {
		log.Fatal("LLAMA_API_KEY environment variable is not set")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	//Execute http request to llama
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get http response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %v", res.Status)
	}

	//Decode http response body into Go struct
	defer res.Body.Close()
	chatRes := &chatResponse{}
	err = json.NewDecoder(res.Body).Decode(chatRes)
	if err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	if len(chatRes.Choices) == 0 {
		log.Fatal("No choices returned from llama")
	}

	//Return generated text from llama
	return chatRes.Choices[0].Message.Content
}

// Set a prompt and other values to create chat request
func createChatRequest(prompt string) *chatRequest {
	return &chatRequest{
		Model: "llama3-70b",
		Messages: []reqMessage{
			reqMessage{Role: "user", Content: prompt},
		},
		Functions: []function{
			function{
				Name:        "Get_English_Exmple_Sentence",
				Description: "Get the English example sentence generated with given words.",
				Parameters: parameters{
					Type: "object",
					Properties: properties{
						Words: words{
							Type:        "string",
							Description: "English vocabulary list, e.g. nonchalant, reckon, appalled",
						},
					},
				},
				Required: []string{"words"},
			},
		},
		Stream:       false,
		FunctionCall: "none",
	}
}

func main() {
	words := [3]string{"nonchalant", "reckon", "appalled"}
	prompt := fmt.Sprintf("Please create an English example sentence using following words: %s, %s, %s",
		words[0], words[1], words[2])

	fmt.Println("")
	fmt.Println("")

	fmt.Println("++++++ Prompt ++++++")
	fmt.Println(prompt)

	fmt.Println("")
	fmt.Println("")

	fmt.Println("++++++ Generated response ++++++")
	fmt.Println(getGeneratedResponse(prompt))

	fmt.Println("")
	fmt.Println("")
}
