package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"resume-review-api/util/resume_ai_env"
	"strings"
	"time"
)

func CreateGPTRequest(messages []Message) (string, error) {
	serverSettings := resume_ai_env.GetSettingsForEnv()

	//var model = os.Getenv("gpt_version")
	var URL = "https://vdartai.openai.azure.com/openai/deployments/" + serverSettings.AzureDeploymentName + "/chat/completions?api-version=2023-05-15"

	jsonData, err := json.Marshal(JSONData{
		Messages: messages,
	})

	if err != nil {
		log.Println("[OpenAI] " + err.Error())
		return "", err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))

	req.Header.Set("api-key", serverSettings.AzureKey)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println("[OpenAI] " + err.Error())
		return "", err
	}

	client := &http.Client{
		Timeout: time.Minute * 5,
	}

	response, err := client.Do(req)

	if err != nil {
		log.Println("[OpenAI] " + err.Error())
		return "", err
	}

	resBody, err := io.ReadAll(response.Body)

	fmt.Println(string(resBody))

	if err != nil {
		log.Println("[OpenAI] " + err.Error())
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
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
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
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}
