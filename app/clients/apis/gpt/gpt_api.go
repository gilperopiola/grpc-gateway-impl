package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
)

/* ———————————————————————————————— — — — GPT API — — — ———————————————————————————————— */

var _ core.GPTAPI = &gptAPI{}

type gptAPI struct {
	httpClient *http.Client
	baseURL    string
	key        string
}

func NewAPI(httpClient *http.Client, apiKey string) core.GPTAPI {
	return &gptAPI{
		httpClient: httpClient, baseURL: "https://api.openai.com/v1", key: apiKey,
	}
}

/* -~-~-~- Endpoints -~-~-~- */

var chatInstructions = apimodels.GPTChatMsg{Role: "user", Content: `You are a highly intelligent and useful AI, showing expertise and excellence in all fields. 
Keep your answers concise and to the point, without repetition. At the beginning of each message, write Title: and a short title for it.`}

func (api *gptAPI) SendRequestToGPT(ctx context.Context, prompt string, prevMsgs ...apimodels.GPTChatMsg) (string, error) {

	gptModel := apimodels.GPT_O1_MINI

	// If no previous messages are provided, we use a default starting message.
	if len(prevMsgs) == 0 {
		prevMsgs = append(prevMsgs, chatInstructions)

		if gptModel == apimodels.GPT_4O || gptModel == apimodels.GPT_4O_MINI {
			prevMsgs[0].Role = "system"
		}
	}

	var url = api.baseURL + "/chat/completions"
	var req = apimodels.GPTChatEndpointRequest{
		Model:    gptModel,
		Messages: append(prevMsgs, apimodels.GPTChatMsg{Role: "user", Content: prompt}),
	}

	status, body, err := utils.POST(ctx, url, &req, api.key, api.httpClient)
	if err != nil {
		return "", logs.LogUnexpected(err)
	}
	logs.LogAPICall(url, status, body)

	var response apimodels.GPTChatEndpointResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshalling gpt chat response: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in gpt chat response")
	}

	return response.Choices[0].Message.Content, nil
}

func (api *gptAPI) SendRequestToDallE(ctx context.Context, prompt string, size pbs.GPTImageSize) (apimodels.GPTImageMsg, error) {
	url := api.baseURL + "/images/generations"
	req := apimodels.GPTImageEndpointRequest{
		Model:  apimodels.DALL_E3,
		Prompt: prompt,
		Size:   api.imageSizeToActualPixels(size),
		N:      1,
	}
	if size == pbs.GPTImageSize_GPT_IMAGE_SIZE_TINY || size == pbs.GPTImageSize_GPT_IMAGE_SIZE_SMALL {
		req.Model = apimodels.DALL_E2
		req.N = 4
	}

	status, body, err := utils.POST(ctx, url, &req, api.key, api.httpClient)
	if err != nil {
		return apimodels.GPTImageMsg{}, logs.LogUnexpected(err)
	}
	logs.LogAPICall(url, status, body)

	var response apimodels.GPTImageEndpointResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return apimodels.GPTImageMsg{}, fmt.Errorf("error unmarshalling dall-e response: %w. request: %+v", err, req)
	}
	if len(response.Data) == 0 {

		// Dall-E-2
		if req.Model == apimodels.DALL_E2 {
			if len(response.ImageURLs) == 0 {
				return apimodels.GPTImageMsg{}, fmt.Errorf("no image URLs in dall-e 2 response")
			}
			return apimodels.GPTImageMsg{URL: response.ImageURLs[0], RevisedPrompt: response.RevisedPrompt}, nil
		}
		return apimodels.GPTImageMsg{}, fmt.Errorf("no data in dall-e 3 response")
	}

	// Dall-E-3
	return apimodels.GPTImageMsg{URL: response.Data[0].URL, RevisedPrompt: response.Data[0].RevisedPrompt}, nil
}

func (api *gptAPI) imageSizeToActualPixels(size pbs.GPTImageSize) string {
	switch size {
	case pbs.GPTImageSize_GPT_IMAGE_SIZE_DEFAULT:
		return "1024x1024"
	case pbs.GPTImageSize_GPT_IMAGE_SIZE_WIDE:
		return "1792x1024"
	case pbs.GPTImageSize_GPT_IMAGE_SIZE_TALL:
		return "1024x1792"
	case pbs.GPTImageSize_GPT_IMAGE_SIZE_SMALL:
		return "512x512"
	case pbs.GPTImageSize_GPT_IMAGE_SIZE_TINY:
		return "256x256"
	default:
		logs.LogUnexpected(fmt.Errorf("invalid dall-eimage size: %v", size))
		return "1024x1024"
	}
}
