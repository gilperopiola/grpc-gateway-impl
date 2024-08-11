package apimodels

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - OpenAI GPT API Models -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type Model string

const (
	GPT35Turbo Model = "gpt-3.5-turbo"
	GPT4oMini  Model = "gpt-4o-mini"
	GPT4Turbo  Model = "gpt-4-turbo"
	GPT4o      Model = "gpt-4o"

	DallE3 Model = "dall-e-3"
)

type (
	Request struct {
		Model    Model     `json:"model"`
		Messages []Message `json:"messages"`
	}

	Response struct {
		ID      string     `json:"id"`
		Choices []Choice   `json:"choices"`
		Created int64      `json:"created"`
		Usage   TokenUsage `json:"usage"`
	}
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type TokenUsage struct {
	InPrompt     int `json:"prompt_tokens"`
	InCompletion int `json:"completion_tokens"`
	InTotal      int `json:"total_tokens"`
}
