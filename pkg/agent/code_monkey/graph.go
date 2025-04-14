package code_monekey

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/JackBekket/GitHelper/pkg/agent/rag/tools"
	graph "github.com/JackBekket/langgraphgo/graph/stategraph"
)

// global var
var Model openai.LLM
var Tools []llms.Tool


// This is the main function for this package
func OneShotRun(prompt string, model openai.LLM, historyState ...llms.MessageContent) string {

	agent_prompt := `For the following task, make plans that can solve the problem step by step. For each plan, indicate \
which external tool together with tool input to retrieve evidence. You can store the evidence into a \
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


	
	agentState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, agent_prompt),
	}

	/*
	initialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Below a current conversation between user and helpful AI assistant. You (assistant) should help user in any task he/she ask you to do."),
	}

	if len(historyState) > 0 { // if there are previouse message state then we first load it into message state
		// Access the first element of the slice
		// ... use the history variable as needed
		initialState = append(initialState, historyState...) // load history as initial state
		initialState = append(
			initialState,
			agentState..., // append agent system prompt
		)
		initialState = append(
			initialState,
			llms.TextParts(llms.ChatMessageTypeHuman, prompt), //append user input (!)
		)
	} else {
		initialState = append(
			initialState,
			agentState...,
		)
		initialState = append(initialState,
			llms.TextParts(llms.ChatMessageTypeHuman, prompt),
		)
	}
		*/

	Tools, _ = tools.GetTools()

	//Tools = tools
	Model = model

	// MAIN WORKFLOW
	workflow := graph.NewStateGraph()

	workflow.AddNode("plan", getPlan)
	workflow.AddNode("tool", tool_execution)
	workflow.AddNode("solve", solve)
	workflow.AddEdge("plan", "tool")
	workflow.AddEdge("solve", END)
	workflow.AddConditionalEdge("tool", _route)
	workflow.SetEntryPoint("plan")

	app, err := workflow.Compile()
	if err != nil {
		log.Printf("error: %v", err)
		return fmt.Sprintf("error :%v", err)
	}

	response, err := app.Invoke(context.Background(), initialState)
	if err != nil {
		log.Printf("error: %v", err)
		return fmt.Sprintf("error :%v", err)
	}

	lastMsg := response[len(response)-1]
	result := lastMsg.Parts[0]
	resultStr := fmt.Sprintf("%v", result)
	return resultStr
}


