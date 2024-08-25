package apimodels

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - OpenAI GPT API Models -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type GPTModel string

const (
	GPT35Turbo GPTModel = "gpt-3.5-turbo"
	GPT4       GPTModel = "gpt-4"
	GPT4Turbo  GPTModel = "gpt-4-turbo"
	GPT4o      GPTModel = "gpt-4o"
	GPT4oMini  GPTModel = "gpt-4o-mini"

	DallE3 GPTModel = "dall-e-3"
)

/* OpenAI API - Request and Response Models */

type (
	GPTChatRequest struct {
		Model    GPTModel     `json:"model"`
		Messages []GPTMessage `json:"messages"`
	}

	GPTChatResponse struct {
		ID      string        `json:"id"`
		Choices []GPTChoice   `json:"choices"`
		Created int64         `json:"created"`
		Usage   GPTTokenUsage `json:"usage"`
	}
)

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTChoice struct {
	Message      GPTMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type GPTTokenUsage struct {
	InPrompt     int `json:"prompt_tokens"`
	InCompletion int `json:"completion_tokens"`
	InTotal      int `json:"total_tokens"`
}
