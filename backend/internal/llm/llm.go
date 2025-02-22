package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

func sendGPTRequest(prompt string) (ChatCompletionResponse, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return ChatCompletionResponse{}, fmt.Errorf("OPENAI_API_KEY not found in environment variables")
	}

	requestBody := ChatCompletionRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model: "gpt-4o-mini",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error marshalling request body: %v", err)
	}

	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ChatCompletionResponse{}, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var response ChatCompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error parsing response JSON: %v", err)
	}

	return response, nil
}

func FetchBullAnalysis(town string) (string, error) {
	prompt := fmt.Sprintf(
		`You are a financial risk management expert analyzing HDB flats in the following town: %s,
		provide a bull thesis for why buyers should consider buying from that area specifically
		list reasons for why the prices might go up and why the area is desirable to live in
		
		keep your response to 300 words`,
		town,
	)

	response, err := sendGPTRequest(prompt)
	if err != nil {
		return "", err
	}
	responseContent := response.Choices[0].Message.Content
	return responseContent, nil
}

func FetchBearAnalysis(town string) (string, error) {
	prompt := fmt.Sprintf(
		`You are a financial risk management expert analyzing HDB flats in the following town: %s,
		provide a bear thesis for why buyers should consider not buying from that area specifically
		list reasons the prices might come down or why the area is not desirable to live in
		
		keep your response to 300 words.`,
		town,
	)

	response, err := sendGPTRequest(prompt)
	if err != nil {
		return "", err
	}
	responseContent := response.Choices[0].Message.Content
	return responseContent, nil
}
