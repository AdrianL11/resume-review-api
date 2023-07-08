package resume

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateGPTRequest(messages []Message) (string, error) {

	var model = "gpt-3.5-turbo-16k"
	var URL = "https://api.openai.com/v1/chat/completions"

	jsonData, err := json.Marshal(JSONData{
		Model:    model,
		Messages: messages,
	})

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))

	req.Header.Set("Authorization", "Bearer "+os.Getenv("openai-key"))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	response, err := client.Do(req)

	if err != nil {
		return "", err
	}

	resBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	openAIResponse := Response{}
	json.Unmarshal(resBody, &openAIResponse)

	return strings.TrimLeft(openAIResponse.Choices[0].Message.Content, "\n"), nil
}

/*-----------------------------------------------------------------------
 * Structures for ChatGPT Send Data
 *-----------------------------------------------------------------------*/

type JSONData struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"'`
	CompletionTokens int `json:"completion_tokens"`
}
