package code_monkey

import (
	"context"
	"fmt"
	"regexp"

	agent "github.com/JackBekket/GitHelper/pkg/agent/rag"
	"github.com/tmc/langchaingo/llms"
)

const PromptGetPlan = `For the following task, make plans that can solve the problem step by step. For each plan, indicate
which external tool together with tool input to retrieve evidence. You can store the evidence into a
variable #E that can be called by later tools. (Plan, #E1, Plan, #E2, Plan, ...)

Tools can be one of the following:
(1) search[json: {"query": "string"}]: Worker that searches results from Duckduckgo. Useful when you need to find short
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
	Task       string                
	PlanString string                 
	Steps      []ReWOOStep               
	Results    map[string]string 
	Result     string                 
	// TODO: add tools to state
}

type ReWOOStep struct {
	Plan      string
	StepName  string
	Tool      string
	ToolInput string
	//Call	  string
	CallMessage	  llms.MessageContent
	Call 	  string
	Result	  string
}

const (
	GraphPlanName  = "plan"
	GraphToolName  = "tool"
	GraphSolveName = "solve"
)

const LLMToolName = "LLM"

var RegexPattern *regexp.Regexp = regexp.MustCompile(`Plan:\s*(.+)\s*(#E\d+)\s*=\s*(\w+)\s*\[([^\]]+)\]`)

func (lc LLMContext) GetPlan(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(ReWOO)
	task := state.Task

	response, err := lc.LLM.GenerateContent(ctx,
		agent.CreateMessageContentHuman(
			fmt.Sprintf(
				"%s\nList of tools:\n%s\n\n%s",
				PromptGetPlan,
				getToolDesc(*lc.Tools),
				task,
			),
		),
	)
	if err != nil {
		return state, err
	}

	result := response.Choices[0].Content
	//matches := RegexPattern.FindAllString(result, -1)

	matches := RegexPattern.FindAllStringSubmatch(result, -1)

	for _, m := range matches {
		state.Steps = append(state.Steps,
			ReWOOStep{
				// m[0] - full match,
				Plan:      m[1],
				StepName:  m[2],
				Tool:      m[3],
				ToolInput: m[4],
			},
		)

	}

	state.PlanString = result

	return state, nil
}

func getToolDesc(tools []llms.Tool) string {
	desc := ""
	tools = append(tools, llms.Tool{
		Function: &llms.FunctionDefinition{
			Name: "LLM",
			Description: `A pretrained LLM like yourself. Useful when you need to act with general
world knowledge and common sense. Prioritize it when you are confident in solving the problem
yourself. Input can be any instruction.`,
		},
	})
	for idx, tool := range tools {
		desc += fmt.Sprintf("(%d) %s[input]: %s\n", idx, tool.Function.Name, tool.Function.Description)
	}
	return desc
}


func getCurrentTask(state ReWOO) int {
	if len(state.Results) == len(state.Steps) {
		return -1
	}
	return len(state.Results)
}
