package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"agent-coach/internal/llm"
	"agent-coach/internal/models"
	"agent-coach/internal/tool"
)

type BaseAgent struct {
	agentType     models.AgentType
	llmRouter     *llm.Router
	toolExecutor  *tool.ToolExecutor
	tools         []models.Tool
	maxIterations int
}

func NewBaseAgent(agentType models.AgentType, router *llm.Router, executor *tool.ToolExecutor, tools []models.Tool) *BaseAgent {
	return &BaseAgent{
		agentType:     agentType,
		llmRouter:     router,
		toolExecutor:  executor,
		tools:         tools,
		maxIterations: 3,
	}
}

func (a *BaseAgent) Type() models.AgentType {
	return a.agentType
}

func (a *BaseAgent) AvailableTools() []models.Tool {
	return a.tools
}

func (a *BaseAgent) BuildContextPrompt(ctx *models.AgentContext) string {
	var sb strings.Builder

	sb.WriteString("\n\n## Current Context\n\n")

	if ctx.Goal != nil {
		sb.WriteString(fmt.Sprintf("**Current Goal**: %s\n", ctx.Goal.Title))
		if ctx.Goal.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", ctx.Goal.Description))
		}
		sb.WriteString(fmt.Sprintf("Status: %s\n", ctx.Goal.Status))
		if ctx.Goal.Context != nil {
			sb.WriteString("\nUser's onboarding answers:\n")
			for k, v := range ctx.Goal.Context {
				sb.WriteString(fmt.Sprintf("- %s: %v\n", k, v))
			}
		}
		sb.WriteString("\n")
	}

	if len(ctx.TodaysTasks) > 0 {
		sb.WriteString("**Today's Tasks**:\n")
		for _, task := range ctx.TodaysTasks {
			status := "⬜"
			if task.Status == "completed" {
				status = "✅"
			} else if task.Status == "in_progress" {
				status = "🔄"
			}
			sb.WriteString(fmt.Sprintf("- %s %s", status, task.Title))
			if task.EstimatedMinutes != nil {
				sb.WriteString(fmt.Sprintf(" (~%d min)", *task.EstimatedMinutes))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	if len(ctx.Tasks) > 0 && len(ctx.TodaysTasks) == 0 {
		sb.WriteString("**Pending Tasks**:\n")
		count := 0
		for _, task := range ctx.Tasks {
			if task.Status == "pending" || task.Status == "in_progress" {
				sb.WriteString(fmt.Sprintf("- %s\n", task.Title))
				count++
				if count >= 5 {
					sb.WriteString(fmt.Sprintf("... and %d more\n", len(ctx.Tasks)-5))
					break
				}
			}
		}
		sb.WriteString("\n")
	}

	// Stats
	sb.WriteString("**Stats**:\n")
	sb.WriteString(fmt.Sprintf("- Tasks completed: %d\n", ctx.TasksCompleted))
	sb.WriteString(fmt.Sprintf("- Current streak: %d days\n", ctx.StreakDays))
	if ctx.RecentStruggles > 0 {
		sb.WriteString(fmt.Sprintf("- Recent struggles: %d\n", ctx.RecentStruggles))
	}

	return sb.String()
}

func (a *BaseAgent) ExecuteWithLLM(ctx context.Context, systemPrompt string, input *AgentInput) (*AgentOutput, error) {
	var messages []llm.Message
	for _, conv := range input.Context.Conversations {
		messages = append(messages, llm.Message{
			Role:    llm.Role(conv.Role),
			Content: conv.Content,
		})
	}

	messages = append(messages, llm.Message{
		Role:    llm.RoleUser,
		Content: input.Message,
	})

	var finalResponse *llm.CompletionResponse

	for iteration := 0; iteration < a.maxIterations; iteration++ {
		resp, err := a.llmRouter.Complete(ctx, &llm.CompletionRequest{
			SystemPrompt: systemPrompt,
			Messages:     messages,
			Tools:        a.tools,
		})
		if err != nil {
			return nil, fmt.Errorf("LLM completion failed: %w", err)
		}

		finalResponse = resp

		if len(resp.ToolCalls) == 0 {
			break
		}

		assistantContent := resp.Content
		messages = append(messages, llm.Message{
			Role:    llm.RoleAssistant,
			Content: assistantContent,
		})

		for _, tc := range resp.ToolCalls {
			result, err := a.toolExecutor.ExecuteTool(ctx, tc.Function.Name, tc.Function.Arguments, input.Context)
			var resultContent string
			if err != nil {
				resultContent = fmt.Sprintf(`{"error": "%s"}`, err.Error())
			} else {
				resultJSON, jsonErr := json.Marshal(result)
				if jsonErr == nil {
					resultContent = string(resultJSON)
				} else {
					resultContent = fmt.Sprintf(`{"status": "success", "result": %v}`, result)
				}
			}

			messages = append(messages, llm.Message{
				Role:       llm.RoleTool,
				Content:    resultContent,
				ToolCallID: tc.ID,
			})
		}
	}

	output := &AgentOutput{
		Response:  finalResponse.Content,
		AgentType: a.agentType,
	}

	return output, nil
}
