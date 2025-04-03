// This will be prototype to superagent (autonomouse agent, which work with memory and have similar functionality to langchain chains.Run method)
package agent

import (
	"log"

	"github.com/JackBekket/GitHelper/internal/database"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// This function fire One-Shot agent without history context
func OnePunch(model, baseURL, localAIToken, prompt string) {
	llm := CreateGenericLLM(model, baseURL, localAIToken)
	call := OneShotRun(prompt, llm)
	log.Println(call)
}

// this function recive previouse history message state and append new user prompt, than run agent
func RunThread(prompt string, model openai.LLM, history ...llms.MessageContent) ([]llms.MessageContent, string) {

	//model := createGenericLLM()
	call := OneShotRun(prompt, model, history...)
	log.Println(call)
	lastResponse := CreateMessageContentAi(call)
	if len(history) > 0 {
		user_msg := CreateMessageContentHuman(prompt)
		state := append(history, user_msg[0])
		state = append(state, lastResponse...)
		return state, call
	} else {
		user_msg := CreateMessageContentHuman(prompt)
		state := user_msg
		state = append(state, lastResponse...)
		return state, call
	}
}

func CreateThread(prompt string, model openai.LLM, collection_name ...string) ([]llms.MessageContent, string) {
	call := OneShotRun(prompt, model)
	log.Println(call)
	lastResponse := CreateMessageContentAi(call)

	user_msg := CreateMessageContentHuman(prompt)
	state := user_msg
	if len(collection_name) >0 {
	for _,cn := range collection_name {
		collectionState := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, "Available Collection Names: " +cn),
		}
		state = append(state,collectionState...)	
	}
	}
		state = append(state, lastResponse...)
		return state, call
}

func CreateMessageContentAi(content string) []llms.MessageContent {
	intialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeAI, content),
	}
	return intialState
}

func createMessageContentSystem(content string) []llms.MessageContent {
	intialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, content),
	}
	return intialState
}

func CreateMessageContentHuman(content string) []llms.MessageContent {
	intialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, content),
	}
	return intialState
}

// refactor
func CreateGenericLLM(model, baseURL, localAIToken string) openai.LLM {
	modelLLM, err := openai.New(
		openai.WithToken(localAIToken),
		openai.WithBaseURL(baseURL),
		openai.WithModel(model),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return *modelLLM
}


func RunNewAgent(aiToken string, model string, baseURL string, prompt string, collection_name ...string) (*database.ChatSessionGraph, string, error) {
	//cb := &ChainCallbackHandler{}

	if baseURL == "" {
		llm, err := openai.New(
			openai.WithToken(aiToken),
			openai.WithModel(model),
			//openai.WithCallback(cb),
		)
		if err != nil {
			return nil, "error", err
		}
		dialogState, outputText := CreateThread(prompt,*llm,collection_name...)
		//last_msg := dialogState[len(dialogState)-1]

		return &database.ChatSessionGraph{
			ConversationBuffer: dialogState,
		}, outputText, nil
	} else {
		llm, err := openai.New(
			openai.WithToken(aiToken),
			openai.WithModel(model),
			openai.WithBaseURL(baseURL),
			openai.WithAPIVersion("v1"),
			//openai.WithCallback(cb),
		)
		if err != nil {
			return nil, "error", err
		}

		dialogState, outputText := RunThread(prompt, *llm)
		return &database.ChatSessionGraph{
			ConversationBuffer: dialogState,
		}, outputText, nil
	}
}

func ContinueAgent(aiToken string, model string, baseURL string, prompt string, state *database.ChatSessionGraph) (*database.ChatSessionGraph, string, error) {
	//cb := &ChainCallbackHandler{}

	if baseURL == "" {
		llm, err := openai.New(
			openai.WithToken(aiToken),
			openai.WithModel(model),
			//openai.WithCallback(cb),
		)
		if err != nil {
			return nil, "error", err
		}
		dialogState, outputText := RunThread(prompt, *llm, state.ConversationBuffer...)

		return &database.ChatSessionGraph{
			ConversationBuffer: dialogState,
		}, outputText, nil
	} else {
		llm, err := openai.New(
			openai.WithToken(aiToken),
			openai.WithModel(model),
			//openai.WithBaseURL("http://localhost:8080"),
			openai.WithBaseURL(baseURL),
			openai.WithAPIVersion("v1"),
			//openai.WithCallback(cb),
		)
		if err != nil {
			return nil, "error", err
		}

		dialogState, outputText := RunThread(prompt, *llm, state.ConversationBuffer...)
		return &database.ChatSessionGraph{
			ConversationBuffer: dialogState,
		}, outputText, nil
	}
}
