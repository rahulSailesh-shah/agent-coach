package models

// Passed to the LLM to describe the tool
type Tool struct {
	Type        string               `json:"type"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Parameters  map[string]ToolParam `json:"parameters"`
}

type ToolParam struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Enum        []string `json:"enum"`
}

// Returned by the LLM to call the tool
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}
