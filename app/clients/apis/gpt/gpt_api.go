package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
)

var _ core.ChatGPTAPI = &ChatGptAPI{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - GPT API -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type ChatGptAPI struct {
	httpClient *http.Client
	baseURL    string
	key        string
}

func NewAPI(httpClient *http.Client, apiKey string) core.ChatGPTAPI {
	return &ChatGptAPI{
		httpClient: httpClient,
		baseURL:    "https://api.openai.com/v1",
		key:        apiKey,
	}
}

/* -~-~-~- Chat Completion Endpoint -~-~-~- */

const endpoint = "/chat/completions"

// Get a text completion from the OpenAI API.
func (api *ChatGptAPI) SendToGPT(ctx context.Context, prompt string, prevMsgs ...apimodels.GPTMessage) (string, error) {

	// We have 1 API Endpoint for both new conversations
	// and existing ones. If no previous messages are
	// provided, we use a default starting message.
	if len(prevMsgs) == 0 {
		prevMsgs = append(prevMsgs, defaultSystemMessage)
	}

	var url = api.baseURL + endpoint
	var req = apimodels.GPTChatRequest{
		Model: apimodels.GPT4Turbo,
		Messages: append(prevMsgs, apimodels.GPTMessage{
			Role:    "user",
			Content: prompt,
		}),
	}

	status, body, err := utils.POST(ctx, url, &req, api.key, api.httpClient)
	if err != nil {
		return "", logs.LogUnexpected(err)
	}

	logs.LogAPICall(url, status, body)

	var response apimodels.GPTChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshalling chat completion response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in chat completion response")
	}

	return response.Choices[0].Message.Content, nil
}

var defaultSystemMessage = apimodels.GPTMessage{
	Role:    "system",
	Content: "You are a highly intelligent and useful AI, showing expertise and excellence in all fields. Keep your answers concise and to the point, without being repetitive. At the beginning of each message, write Title: and a short title for the message.",
}
