package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"agent-coach/internal/llm"
	"agent-coach/internal/models"
)

type IntentClassifier struct {
	llmRouter *llm.Router
}

func NewIntentClassifier(router *llm.Router) *IntentClassifier {
	return &IntentClassifier{
		llmRouter: router,
	}
}

func (c *IntentClassifier) Classify(ctx context.Context, message string, agentCtx *models.AgentContext) (*ClassifiedIntent, error) {
	systemPrompt := c.buildSystemPrompt(agentCtx)
	userMessage := c.buildUserMessage(message, agentCtx)

	resp, err := c.llmRouter.Complete(ctx, &llm.CompletionRequest{
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: userMessage},
		},
	})
	if err != nil {
		return c.fallbackClassification(agentCtx), nil
	}

	result, err := c.parseResponse(resp.Content)
	if err != nil {
		return c.fallbackClassification(agentCtx), nil
	}

	return result, nil
}

func (c *IntentClassifier) buildSystemPrompt(agentCtx *models.AgentContext) string {
	return `You are an intent classifier for an AI coaching application. Your job is to determine what the user wants to do based on their message and the conversation context.

You must choose exactly ONE primary intent for the user's current message.

## Allowed Intent Categories

1. PLANNING
   Choose PLANNING when the user wants to:
   - Create or modify goals, plans, schedules, or study/work roadmaps
   - Set up or adjust a learning path or roadmap
   - Change timelines, priorities, or scope of tasks or goals
   - Ask meta-questions about the plan itself (e.g., "how long will this take?", "what's the plan?", "can we adjust the schedule?")
   - Clarify goals or desired outcomes before working

2. EXECUTION
   Choose EXECUTION when the user wants to:
   - Work on current tasks or get the next task
   - Get help, hints, or explanations for a specific problem or task
   - Continue with the current activity the coach suggested
   - Ask "what should I do now?", "what should I do today?", "what's next?"
   - Confirm readiness to proceed (e.g., "yes", "let's go", "I'm ready", "sounds good") in the context of tasks or execution
   - Provide answers or attempts while actively working on a task

3. EVALUATION
   Choose EVALUATION when the user wants to:
   - Review their progress or performance over time
   - See statistics, streaks, completion rates, or summaries of what they have done
   - Get feedback on how they are doing overall
   - Reflect on what has been working or not working

4. GENERAL
   Choose GENERAL when the user is:
   - Making casual conversation, small talk, or greetings (e.g., "hi", "how are you?")
   - Asking off-topic questions unrelated to their goals, tasks, or progress
   - Saying something unclear, extremely vague, or unrelated to coaching
   - Expressing emotions or talking about motivation without a clear request for a plan, help on a task, or progress review

## IMPORTANT: NO ACCOUNTABILITY INTENT

There is an ACCOUNTABILITY category in the system, but you MUST NOT RETURN IT.

Even if the user is:
- Expressing emotions (frustrated, overwhelmed, unmotivated)
- Requesting a break or saying something is too hard or too easy
You must still classify into PLANNING, EXECUTION, EVALUATION, or GENERAL, based on what they seem to want to do next.

Examples:
- If they say "this plan is too hard, can we simplify it?" → PLANNING (they want to adjust the plan).
- If they say "this task is too hard, I don't know how to do it" → EXECUTION (they want help with a task).
- If they say "I feel burned out, how am I doing lately?" → EVALUATION (they want to review progress).
- If they say "I'm so tired today." with no clear action request → GENERAL.

Never output "ACCOUNTABILITY" as the intent.

## Disambiguation Rules

Use the recent conversation and session state to resolve short or ambiguous messages:

- If the user is responding to a question or confirming something (like "yes", "okay", "sounds good", "let's do it"):
  - Look at what the coach just asked or proposed.
  - If the coach proposed a plan, roadmap, or schedule and the user agrees → PLANNING if the plan is still being constructed or adjusted; EXECUTION once the plan is clearly settled and the user is ready to start.
  - If the coach asked a clarifying question during planning (e.g., about goals, constraints, preferences), and the user answers → usually PLANNING.
  - If the coach was guiding them through a task and the user responds or continues → EXECUTION.

- If the user asks "what should I do now?" or "what should I do today?" → EXECUTION.
- If the user asks to "change the plan", "update my goals", or "rethink the schedule" → PLANNING.
- If the user asks "how am I doing overall?" or "can you show my progress?" → EVALUATION.
- If the user only vents feelings without requesting help with planning, execution, or evaluation → GENERAL.

When in doubt:
- Prefer EXECUTION if the user is currently working on tasks or responding during task work.
- Prefer PLANNING if no tasks exist yet or the conversation is about setting things up.
- Prefer EVALUATION only when the user clearly wants a review or feedback.
- Otherwise, use GENERAL.

## Output Format

Respond with ONLY valid JSON in this exact format (no explanation text, no markdown):

{"intent": "CATEGORY", "confidence": 0.0-1.0, "reason": "brief explanation"}

Where:
- "CATEGORY" is exactly one of: "PLANNING", "EXECUTION", "EVALUATION", "GENERAL"
- "confidence" is a number between 0.0 and 1.0
- "reason" is a short natural-language explanation of why you chose this category.`
}

func (c *IntentClassifier) buildUserMessage(message string, agentCtx *models.AgentContext) string {
	var sb strings.Builder

	sb.WriteString("## Current Context\n\n")

	if agentCtx.Goal != nil {
		sb.WriteString(fmt.Sprintf("**Current Goal**: %s\n", agentCtx.Goal.Title))
	} else {
		sb.WriteString("**Current Goal**: None set yet\n")
	}

	sb.WriteString(fmt.Sprintf("**Tasks**: %d total, %d completed\n", len(agentCtx.Tasks), agentCtx.TasksCompleted))
	sb.WriteString(fmt.Sprintf("**Session State**: %s\n", agentCtx.CurrentState))

	if len(agentCtx.Conversations) > 0 {
		sb.WriteString("\n## Recent Conversation\n\n")
		start := 0
		if len(agentCtx.Conversations) > 5 {
			start = len(agentCtx.Conversations) - 5
		}
		for _, conv := range agentCtx.Conversations[start:] {
			role := "User"
			if conv.Role == "assistant" {
				role = "Coach"
			}
			content := conv.Content
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			sb.WriteString(fmt.Sprintf("**%s**: %s\n\n", role, content))
		}
	}

	sb.WriteString(fmt.Sprintf("## Current User Message\n\n\"%s\"", message))

	return sb.String()
}

func (c *IntentClassifier) parseResponse(content string) (*ClassifiedIntent, error) {
	content = strings.TrimSpace(content)

	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		content = content[start : end+1]
	}

	var result struct {
		Intent     string  `json:"intent"`
		Confidence float64 `json:"confidence"`
		Reason     string  `json:"reason"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	intent := Intent(strings.ToUpper(result.Intent))
	switch intent {
	case IntentPlanning, IntentExecution, IntentEvaluation, IntentAccountability, IntentGeneral:
	default:
		return nil, fmt.Errorf("invalid intent: %s", result.Intent)
	}

	return &ClassifiedIntent{
		Intent:     intent,
		Confidence: result.Confidence,
		Reason:     result.Reason,
	}, nil
}

func (c *IntentClassifier) fallbackClassification(agentCtx *models.AgentContext) *ClassifiedIntent {
	if agentCtx.Goal == nil {
		return &ClassifiedIntent{
			Intent:     IntentPlanning,
			Confidence: 0.6,
			Reason:     "Fallback: No goal set, defaulting to planning",
		}
	}

	if len(agentCtx.Tasks) == 0 {
		return &ClassifiedIntent{
			Intent:     IntentPlanning,
			Confidence: 0.6,
			Reason:     "Fallback: No tasks created yet, continuing planning",
		}
	}

	return &ClassifiedIntent{
		Intent:     IntentExecution,
		Confidence: 0.6,
		Reason:     "Fallback: Goal and tasks exist, defaulting to execution",
	}
}
