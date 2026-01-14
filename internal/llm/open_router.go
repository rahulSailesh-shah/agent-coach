package llm

import (
	"context"
	"encoding/json"

	"agent-coach/internal/models"

	"github.com/google/martian/log"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenRouterProvider struct {
	config *models.LLMProviderConfig
	client openai.Client
}

func NewOpenRouterProvider(config *models.LLMProviderConfig) Provider {
	client := openai.NewClient(
		option.WithAPIKey("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"),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)
	return &OpenRouterProvider{config: config, client: client}
}

func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

func (p *OpenRouterProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	params := p.buildParams(req)
	completion, err := p.client.Chat.Completions.New(ctx, *params)
	if err != nil {
		return nil, err
	}

	toolCalls := make([]ToolCall, len(completion.Choices[0].Message.ToolCalls))
	for i, toolCall := range completion.Choices[0].Message.ToolCalls {
		var arguments map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			log.Errorf("failed to unmarshal tool call arguments: %v", err)
			continue
		}
		toolCalls[i] = ToolCall{
			ID:   toolCall.ID,
			Type: toolCall.Type,
			Function: ToolFunction{
				Name:      toolCall.Function.Name,
				Arguments: arguments,
			},
		}
	}

	return &CompletionResponse{
		Content:      completion.Choices[0].Message.Content,
		ToolCalls:    toolCalls,
		Model:        completion.Model,
		Usage:        int(completion.Usage.TotalTokens),
		FinishReason: completion.Choices[0].FinishReason,
	}, nil
}

func (p *OpenRouterProvider) IsAvailable() bool {
	return p.config.IsActive
}

func (p *OpenRouterProvider) buildParams(req *CompletionRequest) *openai.ChatCompletionNewParams {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(req.Messages)+1)
	messages[0] = openai.SystemMessage(req.SystemPrompt)
	for _, message := range req.Messages {
		switch message.Role {
		case RoleUser:
			messages = append(messages, openai.UserMessage(message.Content))
		case RoleAssistant:
			messages = append(messages, openai.AssistantMessage(message.Content))
		case RoleTool:
			messages = append(messages, openai.ToolMessage(message.Content, message.ToolCallID))
		default:
			continue
		}
	}

	tools := make([]openai.ChatCompletionToolUnionParam, len(req.Tools))
	for _, tool := range req.Tools {
		switch tool.Type {
		case "function":
			tools = append(tools, buildFunctionTool(tool))
		default:
			continue
		}
	}

	model := req.Model
	if model == "" {
		model = p.config.DefaultModel
	}
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Tools:    tools,
		Model:    openai.ChatModel(model),
	}

	return &params

}

func buildFunctionTool(tool Tool) openai.ChatCompletionToolUnionParam {
	required := make([]string, 0)
	properties := make(map[string]any)
	for name, param := range tool.Parameters {
		properties[name] = map[string]any{
			"type":        param.Type,
			"description": param.Description,
			"enum":        param.Enum,
		}
		if param.Required {
			required = append(required, name)
		}
	}
	parameters := openai.FunctionParameters{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        tool.Name,
		Description: openai.String(tool.Description),
		Parameters:  parameters,
	})
}
