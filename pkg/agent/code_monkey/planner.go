package code_monkey

import (
	"context"
	"fmt"
	"regexp"

	agent "github.com/JackBekket/GitHelper/pkg/agent/rag"
	"github.com/tmc/langchaingo/llms"
)

const AgentPrompt = `For the following task, make plans that can solve the problem step by step. For each plan, indicate
which external tool together with tool input to retrieve evidence. You can store the evidence into a
variable #E that can be called by later tools. (Plan, #E1, Plan, #E2, Plan, ...)

Tools can be one of the following:
(1) Google[input]: Worker that searches results from Google. Useful when you need to find short
and succinct answers about a specific topic. The input should be a search query.
(2) LLM[input]: A pretrained LLM like yourself. Useful when you need to act with general
world knowledge and common sense. Prioritize it when you are confident in solving the problem
yourself. Input can be any instruction.

For example,
Task: Thomas, Toby, and Rebecca worked a total of 157 hours in one week. Thomas worked x
hours. Toby worked 10 hours less than twice what Thomas worked, and Rebecca worked 8 hours
less than Toby. How many hours did Rebecca work?
Plan: Given Thomas worked x hours, translate the problem into algebraic expressions and solve
with Wolfram Alpha. #E1 = WolframAlpha[Solve x + (2x − 10) + ((2x − 10) − 8) = 157]
Plan: Find out the number of hours Thomas worked. #E2 = LLM[What is x, given #E1]
Plan: Calculate the number of hours Rebecca worked. #E3 = Calculator[(2 ∗ #E2 − 10) − 8]

Begin! 
Describe your plans with rich details. Each Plan should be followed by only one #E.

Task: `

type ReWOO struct {
	Task       string                 `json:"task"`
	PlanString string                 `json:"plan_string"`
	Steps      []string               `json:"steps"`
	Results    map[string]interface{} `json:"results"`
	Result     string                 `json:"result"`
}

var RegexPattern *regexp.Regexp = regexp.MustCompile(`Plan:\s*(.+)\s*(#E\d+)\s*=\s*(\w+)\s*\[([^\]]+)\]`)

func (lc LLMContext) GetPlan(ctx context.Context, state interface{}) (interface{}, error) {
	rwState := state.(ReWOO)
	task := rwState.Task

	response, err := lc.LLM.GenerateContent(ctx,
		agent.CreateMessageContentHuman(
			fmt.Sprintf(
				"%s\nList of tools:\n%s\n\n%s",
				AgentPrompt,
				getToolDesc(*lc.Tools),
				task,
			),
		),
	)
	if err != nil {
		return state, err
	}

	result := response.Choices[0].Content
	matches := RegexPattern.FindAllString(result, -1)

	rwState.Steps = matches
	rwState.PlanString = result

	return rwState, nil
}

func getToolDesc(tools []llms.Tool) string {
	desc := ""
	for idx, tool := range tools {
		desc += fmt.Sprintf("(%d) %s[input]: %s\n", idx, tool.Function.Name, tool.Function.Description)
	}
	return desc
}


func _getCurrentTask(state ReWOO) int {
    if state.Results == nil {
        return 1
    }
    if len(state.Results) == len(state.Steps) {
        return 0
    } else {
        return int(len(state.Results) + 1)
    }
}
