package tool

import "agent-coach/internal/models"

var ToolCreateMilestone = models.Tool{
	Name:        "create_milestone",
	Description: "Create a milestone (major goal checkpoint) with a target week range",
	Parameters: map[string]models.ToolParam{
		"title":       {Type: "string", Description: "Milestone title", Required: true},
		"description": {Type: "string", Description: "What this milestone covers", Required: true},
		"week_start":  {Type: "integer", Description: "Starting week number", Required: true},
		"week_end":    {Type: "integer", Description: "Ending week number", Required: true},
	},
}

var ToolCreateTask = models.Tool{
	Name:        "create_task",
	Description: "Create a specific task for the user to complete",
	Parameters: map[string]models.ToolParam{
		"title":             {Type: "string", Description: "Task title", Required: true},
		"description":       {Type: "string", Description: "Task description", Required: false},
		"due_date":          {Type: "string", Description: "Due date in YYYY-MM-DD format", Required: false},
		"estimated_minutes": {Type: "integer", Description: "Estimated time in minutes", Required: false},
		"difficulty":        {Type: "integer", Description: "Difficulty 1-5", Required: false},
		"priority":          {Type: "integer", Description: "Priority (higher = more important)", Required: false},
	},
}

var ToolSuggestResources = models.Tool{
	Name:        "suggest_resources",
	Description: "Suggest learning resources to the user",
	Parameters: map[string]models.ToolParam{
		"resources": {Type: "array", Description: "List of resource objects with title, type, url", Required: true},
	},
}

var ToolAskClarifyingQuestion = models.Tool{
	Name:        "ask_clarifying_question",
	Description: "Ask the user a clarifying question to better understand their needs",
	Parameters: map[string]models.ToolParam{
		"question": {Type: "string", Description: "The question to ask", Required: true},
	},
}
