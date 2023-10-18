package chat

// Message represents a chat message
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// CompletionChat represents completion chat settings.
type CompletionChat struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens uint      `json:"max_tokens"`
	// ... other fields we don't care about
}

// Choice represents a single choice generated by completion apiutil.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Completion represents a response body returned by chat completion apiutil.
type Completion struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
	// ... other fields we don't care about
}
