package code_monkey

import (
	"context"
	"fmt"
	"strings"

	agent "github.com/JackBekket/GitHelper/pkg/agent/rag"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
)

func (lc LLMContext) ToolExecution(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(ReWOO)

	step := state.Steps[getCurrentTask(state)]

	for stepName, result := range state.Results {
		step.ToolInput = strings.ReplaceAll(step.ToolInput, stepName, result)
	}

	toolchain, err := InitializeChain()
	if err != nil {
		log.Printf("Can't compile toolchain")
		return nil, err
	}

	prompt := step.ToolInput
	options := []llms.CallOption{}
	content := ""
	if step.Tool != LLMToolName {
		prompt = fmt.Sprintf(
			"Use tool %s to process the task.\nTask: %s",	// no use prompt
			step.Tool,
			prompt,
		)
		
		result, err := toolchain.Invoke(ctx,step)
		if err != nil {
			return nil, err
		}
		r := result.(ReWOOStep)
		content = r.Result
	} else {
		response, err := lc.LLM.GenerateContent(ctx,
			agent.CreateMessageContentHuman(
				prompt,
			),
			options...,
		)
		if err != nil {
			return state, err
		}
		content = response.Choices[0].Content
	}

	if len(state.Results) == 0 {
		state.Results = map[string]string{}
	}

	state.Results[step.StepName] = content
	return state, nil
}


