package tool

import "agent-coach/internal/models"

var ToolPresentTask = models.Tool{
	Name:        "present_task",
	Description: "Present a task to the user with context and guidance",
	Parameters: map[string]models.ToolParam{
		"task_id":  {Type: "string", Description: "ID of the task to present", Required: true},
		"approach": {Type: "string", Description: "Suggested approach for the task", Required: false},
		"tips":     {Type: "array", Description: "List of tips for completing the task", Required: false},
	},
}

var ToolProvideHint = models.Tool{
	Name:        "provide_hint",
	Description: "Provide a hint to help the user with their current task",
	Parameters: map[string]models.ToolParam{
		"hint_level": {Type: "integer", Description: "Hint level 1-3 (1=subtle, 3=explicit)", Required: true},
		"hint":       {Type: "string", Description: "The hint content", Required: true},
	},
}

var ToolMarkComplete = models.Tool{
	Name:        "mark_complete",
	Description: "Mark a task as completed",
	Parameters: map[string]models.ToolParam{
		"task_id":        {Type: "string", Description: "ID of the task", Required: true},
		"actual_minutes": {Type: "integer", Description: "Actual time taken in minutes", Required: false},
	},
}

var ToolLogStruggle = models.Tool{
	Name:        "log_struggle",
	Description: "Log that the user struggled with a task",
	Parameters: map[string]models.ToolParam{
		"task_id": {Type: "string", Description: "ID of the task", Required: true},
		"notes":   {Type: "string", Description: "Notes about the struggle", Required: true},
	},
}
