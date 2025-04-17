package code_monkey

import (
	"context"
	"os"

	agent "github.com/JackBekket/GitHelper/pkg/agent/rag"
	"github.com/JackBekket/GitHelper/pkg/agent/rag/tools"
	graph "github.com/JackBekket/langgraphgo/graph/stategraph"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)


func initializeModel() (openai.LLM){

	api_token := os.Getenv("API_TOKEN")
	ai_url := os.Getenv("AI_URL")
	model := os.Getenv("MODEL")
	llm := agent.CreateGenericLLM(model,ai_url,api_token)
	return llm
}


func InitializeChain() (*graph.Runnable,error){

		// MAIN WORKFLOW
		workflow := graph.NewStateGraph()

		workflow.AddNode("generate_call", generateCall)
		workflow.AddNode("semanticSearch",semanticSearch)

		workflow.AddEdge("generate_call","semanticSearch")
		workflow.AddEdge("semanticSearch","END")
	
		app, err := workflow.Compile()
		if err != nil {
			log.Printf("error: %v", err)
			return nil,err
		}
		return app,nil
	
		/*
		response, err := app.Invoke(context.Background(), initialState)
		if err != nil {
			log.Printf("error: %v", err)
			return fmt.Sprintf("error :%v", err)
		}
	
		lastMsg := response[len(response)-1]
		result := lastMsg.Parts[0]
		resultStr := fmt.Sprintf("%v", result)
		return resultStr
		*/

}



// This function is performing similarity search in our db vectorstore.
func semanticSearch(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(ReWOOStep)
	semanticSearchTool := tools.SemanticSearchTool{}
	call := state.Call
	input := agent.CreateMessageContentHuman(call)
	res, err := semanticSearchTool.Execute(ctx, input)
	if err != nil {
		// Handle the error
		return nil, err
	}
	content := ""
  for _, msg := range res {
    if msg.Role == llms.ChatMessageTypeTool {
      for _, part := range msg.Parts {
        if toolCallResponse, ok := part.(llms.ToolCallResponse); ok {
          content += toolCallResponse.Content
        }
      }
    }
  }
  	state.Result = content
	return state, nil
}


func generateCall(ctx context.Context, s interface{}) (interface{},error){
	state := s.(ReWOOStep)
	prompt := state.ToolInput
	model := initializeModel()
	tools, err := tools.GetTools()
	if err != nil {
		log.Printf("error getting tools", err)
	}
	msg := agent.CreateMessageContentHuman(prompt)
	response, err := model.GenerateContent(ctx, msg, llms.WithTools(tools)) // AI call tool function.. in this step it just put call in messages stack
	if err != nil {
		log.Printf("error generating tool call", err)
		return state,err
	}
	result := llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content)
	for _, part := range result.Parts {
		toolCall, ok := part.(llms.ToolCall)

		if ok && toolCall.FunctionCall.Name == "semanticSearch" {
			log.Printf("agent should use SemanticSearch (embeddings similarity search aka DocumentsSearch)")
			
		}
	}
	state.Call = response.Choices[0].Content
	return state,nil
}

