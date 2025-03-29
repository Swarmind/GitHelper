// This will be prototype to superagent (autonomouse agent, which work with memory and have similar functionality to langchain chains.Run method)
package agent

import (
	"log"

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
			llms.TextParts(llms.ChatMessageTypeSystem, "Collection Name: " +cn),
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
