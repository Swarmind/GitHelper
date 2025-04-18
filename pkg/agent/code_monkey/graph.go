package code_monkey

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/JackBekket/GitHelper/pkg/agent/rag/tools"
	graph "github.com/JackBekket/langgraphgo/graph/stategraph"
)

type LLMContext struct {
	LLM           *openai.LLM
	Tools         *[]llms.Tool
	ToolsExecutor *tools.ToolsExectutor			// ??????
}

func (lc LLMContext) OneShotRun(ctx context.Context, prompt string) (string, error) {

	workflowGraph := graph.NewStateGraph()

	workflowGraph.AddNode("plan", lc.GetPlan)
	workflowGraph.AddNode("tool", lc.ToolExecution)
	workflowGraph.AddNode("solve", lc.Solve)
	workflowGraph.SetEntryPoint("plan")
	workflowGraph.AddEdge("plan", "tool")
	workflowGraph.AddConditionalEdge("tool", Route)
	workflowGraph.AddEdge("solve", graph.END)
	
	

	app, err := workflowGraph.Compile()
	if err != nil {
		return "", err
	}

	state, err := app.Invoke(ctx, ReWOO{
		Task: prompt,
	})
	if err != nil {
		return "", err
	}

	return state.(ReWOO).Result, nil
}