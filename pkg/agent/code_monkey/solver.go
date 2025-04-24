package code_monkey

import (
	"context"
	"fmt"
	"strings"

	agent "github.com/JackBekket/GitHelper/pkg/agent/rag"
)

const PromptSolver = `Solve the following task or problem. To solve the problem, we have made step-by-step Plan and \
retrieved corresponding Evidence to each Plan. Use them with caution since long evidence might \
contain irrelevant information.

%s

Now solve the question or task according to provided Evidence above. Respond with the answer
directly with no extra words.

Task: %s
Response:`



func (lc LLMContext) Solve(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(ReWOO)

	plan := ""
	for _, step := range state.Steps {
		for stepName, result := range state.Results {
			step.ToolInput = strings.ReplaceAll(step.ToolInput, stepName, result)
			step.StepName = strings.ReplaceAll(step.StepName, stepName, result)
		}
		plan += fmt.Sprintf("Plan: %s\n%s = %s[%s]\n", step.Plan, step.StepName, step.Tool, step.ToolInput)
	}
	response, err := lc.LLM.GenerateContent(ctx,
		agent.CreateMessageContentHuman(
			fmt.Sprintf(PromptSolver, plan, state.Task),
		),
	)
	if err != nil {
		return state, err
	}

	state.Result = response.Choices[0].Content

	return state, nil
}

func Route(ctx context.Context, state interface{}) string {
	if getCurrentTask(state.(ReWOO)) == -1 {
		return GraphSolveName
	} else {
		return GraphToolName
	}
}