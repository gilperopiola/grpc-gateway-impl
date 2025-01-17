package apimodels

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - OpenAI GPT API Models -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type GPTs string

const (
	GPT_O1_PREVIEW GPTs = "o1-preview"
	GPT_O1_MINI    GPTs = "o1-mini"
	GPT_4O         GPTs = "gpt-4o"
	GPT_4O_MINI    GPTs = "gpt-4o-mini"

	DALL_E2 GPTs = "dall-e-2"
	DALL_E3 GPTs = "dall-e-3"
)

// Data Models for the Chat Completions API
type (
	GPTChatEndpointRequest struct {
		Model    GPTs         `json:"model"`
		Messages []GPTChatMsg `json:"messages"`
	}

	GPTChatEndpointResponse struct {
		ID      string/*                */ `json:"id"`
		Choices []struct {
			Message      GPTChatMsg/*   */ `json:"message"`
			FinishReason string/*       */ `json:"finish_reason"`
		}/*                            	*/ `json:"choices"`
		Usage struct {
			InPrompt     int/*          */ `json:"prompt_tokens"`
			InCompletion int/*          */ `json:"completion_tokens"`
			InTotal      int/*          */ `json:"total_tokens"`
		}/*                            	*/ `json:"usage"`
		Created int64/*                 */ `json:"created"`
	}

	GPTChatMsg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

// Data Models for the Images API
type (
	GPTImageEndpointRequest struct {
		Model  GPTs   `json:"model"`
		Prompt string `json:"prompt"`
		N      int    `json:"n"`
		Size   string `json:"size"`
	}

	GPTImageEndpointResponse struct {
		Data    []GPTImageMsg `json:"data"`
		Created int64         `json:"created"`

		// Dall-E 2 has a different response format:
		RevisedPrompt string   `json:"revised_prompt"`
		ImageURLs     []string `json:"image_urls"`
	}

	GPTImageMsg struct {
		URL           string `json:"url"`
		RevisedPrompt string `json:"revised_prompt"`
	}
)
