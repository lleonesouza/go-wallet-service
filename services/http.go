package services

import (
	"bff-answerfy/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Http struct {
	env *config.Envs
}

type Messages struct {
	role    string
	content string
}

type ChatInput struct {
	model       string
	messages    []Messages
	temperature int
	max_tokens  int
}

type ChatResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func readBytes(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func Unmarshal(body []byte, data any) error {
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	return nil
}

func Marshal(data map[string]string) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (h *Http) OpenAIRequest(text string) (*ChatResponse, error) {
	url := "https://api.openai.com/v1/completions"
	data := map[string]interface{}{
		"model":       "text-davinci-003",
		"prompt":      text,
		"max_tokens":  2048,
		"temperature": 0,
	}

	// Set the API endpoint URL and request body
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request with the specified method, URL, and body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set the API authentication header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.env.OPENAI_API_KEY))
	req.Header.Set("OpenAI-Organization", "org-hYHpOnFG44qHpVx2BTiLp5dA")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body into a custom struct
	var chatResponse ChatResponse
	err = json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		return nil, err
	}

	return &chatResponse, nil
}

func (h *Http) POST(url string, body []byte) (*http.Response, error) {
	// Create a new POST request with the JSON body
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set the request header to JSON
	request.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}

func (h *Http) GET(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (h *Http) PUT(url string, body []byte) (*http.Response, error) {
	// Create a new POST request with the JSON body
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set the request header to JSON
	request.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}

func (h *Http) DELETE(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp, nil
}

// curl https://api.openai.com/v1/chat/completions \
// -H 'Content-Type: application/json' \
// -H 'Authorization: Bearer x' \
// -H 'OpenAI-Organization: x' \
// -d '{
//   "model": "gpt-3.5-turbo",
//   "messages": [{"role": "user", "content": "Hello!"}],
// "max_tokens": 7,
// "temperature": 0
// }'
