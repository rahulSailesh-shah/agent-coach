package tool

import "agent-coach/internal/models"

var ToolRatePerformance = models.Tool{
	Name:        "rate_performance",
	Description: "Rate the user's performance on a task",
	Parameters: map[string]models.ToolParam{
		"task_id":  {Type: "string", Description: "ID of the task", Required: true},
		"rating":   {Type: "integer", Description: "Performance rating 1-5", Required: true},
		"feedback": {Type: "string", Description: "Feedback for the user", Required: true},
	},
}

var ToolIdentifyWeakness = models.Tool{
	Name:        "identify_weakness",
	Description: "Identify an area where the user needs improvement",
	Parameters: map[string]models.ToolParam{
		"area":       {Type: "string", Description: "The weakness area", Required: true},
		"evidence":   {Type: "string", Description: "Evidence supporting this observation", Required: true},
		"suggestion": {Type: "string", Description: "Suggested action to improve", Required: true},
	},
}

var ToolSuggestReview = models.Tool{
	Name:        "suggest_review",
	Description: "Suggest topics or tasks for review",
	Parameters: map[string]models.ToolParam{
		"topics": {Type: "array", Description: "Topics to review", Required: true},
		"reason": {Type: "string", Description: "Why review is needed", Required: true},
	},
}
