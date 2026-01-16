package tool

import "agent-coach/internal/models"

var ToolSendCheckIn = models.Tool{
	Name:        "send_checkin",
	Description: "Send a check-in message to the user",
	Parameters: map[string]models.ToolParam{
		"message":   {Type: "string", Description: "Check-in message", Required: true},
		"questions": {Type: "array", Description: "Questions to ask", Required: false},
	},
}

var ToolRecordCheckIn = models.Tool{
	Name:        "record_checkin",
	Description: "Record the user's check-in response",
	Parameters: map[string]models.ToolParam{
		"summary":     {Type: "string", Description: "Summary of the check-in", Required: true},
		"struggles":   {Type: "array", Description: "Things the user struggled with", Required: false},
		"wins":        {Type: "array", Description: "Things that went well", Required: false},
		"mood_rating": {Type: "integer", Description: "Mood rating 1-5", Required: false},
	},
}

var ToolAdjustPace = models.Tool{
	Name:        "adjust_pace",
	Description: "Adjust the pace of the user's plan",
	Parameters: map[string]models.ToolParam{
		"adjustment": {Type: "string", Description: "Type of adjustment: slower, faster, maintain", Required: true, Enum: []string{"slower", "faster", "maintain"}},
		"reason":     {Type: "string", Description: "Reason for adjustment", Required: true},
	},
}

var ToolCelebrateWin = models.Tool{
	Name:        "celebrate_win",
	Description: "Celebrate a user achievement",
	Parameters: map[string]models.ToolParam{
		"achievement": {Type: "string", Description: "What to celebrate", Required: true},
		"message":     {Type: "string", Description: "Celebratory message", Required: true},
	},
}
