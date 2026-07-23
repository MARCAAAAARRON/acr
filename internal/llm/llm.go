package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const baseURL = "http://localhost:1234/v1/chat/completions"

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

// Ask sends prompt to the given model and returns its reply text.
func Ask(model string, prompt string) (string, error) {
	// 1. build the request struct
	chatReq := chatRequest{
		Model: model,
		Messages: []message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 2. turn it into JSON bytes — must happen before we can POST it
	data, err := json.Marshal(chatReq)
	if err != nil {
		return "", err
	}

	// 3. send it — this is the first line allowed to use `data`
	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 4. decode the JSON response into our struct
	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// 5. hand back the actual reply text
	return result.Choices[0].Message.Content, nil
}