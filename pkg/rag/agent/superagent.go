// This will be prototype to superagent (autonomouse agent, which work with memory and have similar functionality to langchain chains.Run method)
package agent

import (
	"fmt"
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

func CreateThread(prompt string, model openai.LLM, collectionName ...string) ([]llms.MessageContent, string) {
	historyState := []llms.MessageContent{}
	if len(collectionName) > 0 {
		msg := "Available Collection Names:\n"
		for _, cn := range collectionName {
			msg += fmt.Sprintf("%s\n", cn)
		}
		historyState = append(historyState, llms.MessageContent{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextPart(msg),
			},
		})
	}

	call := OneShotRun(prompt, model, historyState...)
	log.Println(call)

	promptRequest := CreateMessageContentHuman(prompt)
	lastResponse := CreateMessageContentAi(call)

	historyState = append(historyState, promptRequest...)
	historyState = append(historyState, lastResponse...)

	return historyState, call
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
		dialogState, outputText := CreateThread(prompt, *llm, collection_name...)
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

		dialogState, outputText := CreateThread(prompt, *llm, collection_name...)
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
