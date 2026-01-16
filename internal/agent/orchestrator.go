package agent

import (
	"agent-coach/internal/tool"
	"context"
	"fmt"
	"log"

	"agent-coach/internal/llm"
	"agent-coach/internal/models"
	"agent-coach/internal/storage"
)

type Orchestrator struct {
	db           *storage.DB
	llmRouter    *llm.Router
	classifier   *IntentClassifier
	toolExecutor *tool.ToolExecutor

	plannerAgent  *PlannerAgent
	executorAgent *ExecutorAgent

	goalRepo *storage.GoalRepository
	taskRepo *storage.TaskRepository
	convRepo *storage.ConversationRepository
}

func NewOrchestrator(db *storage.DB, router *llm.Router) *Orchestrator {
	toolExecutor := tool.NewToolExecutor(db)

	plannerAgent := NewPlannerAgent(toolExecutor)
	executorAgent := NewExecutorAgent(toolExecutor)

	plannerAgent.BaseAgent.llmRouter = router
	executorAgent.BaseAgent.llmRouter = router

	return &Orchestrator{
		db:            db,
		llmRouter:     router,
		classifier:    NewIntentClassifier(router),
		toolExecutor:  toolExecutor,
		plannerAgent:  plannerAgent,
		executorAgent: executorAgent,
		goalRepo:      storage.NewGoalRepository(db),
		taskRepo:      storage.NewTaskRepository(db),
		convRepo:      storage.NewConversationRepository(db),
	}
}

func (o *Orchestrator) ProcessMessage(ctx context.Context, message string, goalID string) (*AgentOutput, error) {
	log.Printf("[Orchestrator] Processing message: %q (goalID=%s)", message, goalID)

	// 1. Build context
	agentCtx, err := o.buildContext(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to build context: %w", err)
	}
	log.Printf("[Orchestrator] Context built: state=%s, goal=%v, tasks=%d",
		agentCtx.CurrentState,
		agentCtx.Goal != nil,
		len(agentCtx.Tasks))

	// 2. Classify intent
	classified, err := o.classifier.Classify(ctx, message, agentCtx)
	if err != nil {
		log.Printf("[Orchestrator] Classification error: %v, defaulting to executor", err)
		classified = &ClassifiedIntent{
			Intent:     IntentExecution,
			Confidence: 0.5,
			Reason:     "Classification failed, defaulting to executor",
		}
	}
	log.Printf("[Orchestrator] Classified intent: %s (confidence=%.2f, reason=%s)",
		classified.Intent, classified.Confidence, classified.Reason)

	// 3. Route to appropriate agent
	agent := o.routeToAgent(classified.Intent)
	log.Printf("[Orchestrator] Routing to agent: %s", agent.Type())

	// 4. Execute agent
	input := &AgentInput{
		Message:   message,
		Context:   agentCtx,
		GoalID:    goalID,
		SessionID: fmt.Sprintf("%s", goalID),
	}

	output, err := agent.Execute(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
	}

	// 5. Save conversation
	o.saveConversation(ctx, goalID, message, output)

	return output, nil
}

func (o *Orchestrator) buildContext(ctx context.Context, goalID string) (*models.AgentContext, error) {
	agentCtx := &models.AgentContext{}

	goal, err := o.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		return nil, err
	}
	agentCtx.Goal = goal

	if agentCtx.Goal != nil {
		tasks, err := o.taskRepo.GetByGoalID(ctx, agentCtx.Goal.ID)
		if err == nil {
			agentCtx.Tasks = tasks
			for _, task := range tasks {
				if task.Status == "completed" {
					agentCtx.TasksCompleted++
				}
				if task.StruggleNotes != "" {
					agentCtx.RecentStruggles++
				}
			}
		}
		todaysTasks, err := o.taskRepo.GetDueToday(ctx)
		if err == nil {
			for _, task := range todaysTasks {
				if task.GoalID == agentCtx.Goal.ID {
					agentCtx.TodaysTasks = append(agentCtx.TodaysTasks, task)
				}
			}
		}
		if len(agentCtx.Tasks) == 0 {
			agentCtx.CurrentState = models.StatePlanning
		}
	} else {
		agentCtx.CurrentState = models.StateGoalSetting
	}

	convs, err := o.convRepo.GetByGoalID(ctx, goalID, 10)
	if err == nil {
		agentCtx.Conversations = convs
	}

	return agentCtx, nil
}

func (o *Orchestrator) routeToAgent(intent Intent) Agent {
	switch intent {
	case IntentPlanning:
		return o.plannerAgent
	case IntentExecution:
		return o.executorAgent
	default:
		return o.executorAgent
	}
}

func (o *Orchestrator) saveConversation(ctx context.Context, goalID string, userMessage string, output *AgentOutput) {
	userConv := &models.Conversation{
		GoalID:    &goalID,
		SessionID: goalID,
		Role:      models.RoleUser,
		Content:   userMessage,
	}
	o.convRepo.Create(ctx, userConv)

	agentConv := &models.Conversation{
		GoalID:    &goalID,
		SessionID: goalID,
		Role:      models.RoleAssistant,
		Content:   output.Response,
		AgentType: output.AgentType,
	}
	o.convRepo.Create(ctx, agentConv)
}

func (o *Orchestrator) GetAgentForIntent(intent Intent) models.AgentType {
	switch intent {
	case IntentPlanning:
		return models.AgentTypePlanner
	case IntentExecution:
		return models.AgentTypeExecutor
	case IntentEvaluation:
		return models.AgentTypeEvaluator
	case IntentAccountability:
		return models.AgentTypeAccountability
	default:
		return models.AgentTypeExecutor
	}
}

func (o *Orchestrator) StartOnboarding(ctx context.Context, goalID string, goalTitle string) (*AgentOutput, error) {
	message := fmt.Sprintf("I want to set a new goal: %s. Help me get started.", goalTitle)
	return o.ProcessMessage(ctx, message, goalID)
}
