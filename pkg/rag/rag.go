package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/JackBekket/hellper/lib/embeddings"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores"
)

// main function for retrieval-augmented generation  (old one)
func Rag(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error) {
	base_url := ai_url

	// Create an embeddings client using the specified API and embedding model
	llm, err := openai.New(
		openai.WithBaseURL(base_url),
		openai.WithAPIVersion("v1"),
		openai.WithToken(api_token),
		openai.WithModel("tiger-gemma-9b-v1-i1"),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	)
	if err != nil {
		return "", err
	}

	//🤕🤕🤕
	searchResults, err := embeddings.SemanticSearch(question, numOfResults, store)
	if err != nil {
		return "", err
	}

	contextBuilder := strings.Builder{}
	for _, doc := range searchResults {
		contextBuilder.WriteString(doc.PageContent)
		contextBuilder.WriteString("\n")
	}
	contexts := contextBuilder.String()

	fullPrompt := fmt.Sprintf("Context: %s\n\nQuestion: %s", contexts, question)

	result, err = chains.Run(
		context.Background(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, numOfResults),
		),
		fullPrompt,
		chains.WithMaxTokens(8192),
	)
	if err != nil {
		return "", err
	}

	fmt.Println("====final answer====\n", result)

	return result, nil
}

func RagWithFilteres(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore)  (result string, err error) {
		//base_url := os.Getenv("AI_BASEURL")
		base_url := ai_url

		// Create an embeddings client using the specified API and embedding model
		llm, err := openai.New(
			openai.WithBaseURL(base_url),
			openai.WithAPIVersion("v1"),
			openai.WithToken(api_token),
			openai.WithModel("tiger-gemma-9b-v1-i1"),
			openai.WithEmbeddingModel("text-embedding-ada-002"),
		)
		if err != nil {
			return "", err
		}
		

		// filteres
		filters := map[string]any{
			"type": "doc",
			//"author": "John Doe",
		}
		

		/*	If we need additional filters 
		options := vectorstores.Options{
			Filters: filters,
		}
		*/

		option := vectorstores.WithFilters(filters)
		  
		searchResults, err := embeddings.SemanticSearch(question, numOfResults, store, option)

		if err != nil {
			return "", err
		}
	
		contextBuilder := strings.Builder{}
		for _, doc := range searchResults {
			contextBuilder.WriteString(doc.PageContent)
			contextBuilder.WriteString("\n")
		}
		contexts := contextBuilder.String()
	
		fullPrompt := fmt.Sprintf("Context: %s\n\nQuestion: %s", contexts, question)
	
		result, err = chains.Run(
			context.Background(),
			chains.NewRetrievalQAFromLLM(
				llm,
				vectorstores.ToRetriever(store, numOfResults),
			),
			fullPrompt,
			chains.WithMaxTokens(8192),
		)
		if err != nil {
			return "", err
		}
	
		fmt.Println("====final answer====\n", result)
	
		return result, nil

}


/*
 Retrieval - Augmented generation, using stuffed QA (prompt is transferred with all page conent it was found)
 Creates a chain that takes input documents and a question.
Combines all documents into a single prompt for the LLM.
Suitable for a small number of documents.
*/
func StuffedQA_Rag(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error) {
		//base_url := os.Getenv("AI_BASEURL")
		base_url := ai_url

		// Create an embeddings client using the specified API and embedding model
		llm, err := openai.New(
			openai.WithBaseURL(base_url),
			openai.WithAPIVersion("v1"),
			openai.WithToken(api_token),
			openai.WithModel("tiger-gemma-9b-v1-i1"),
			openai.WithEmbeddingModel("text-embedding-ada-002"),
		)
		if err != nil {
			return "", err
		}
	
		//docs we found in this store
		searchResults, err := embeddings.SemanticSearch(question, numOfResults, store)
		if err != nil {
			return "", err
		}

		  
	// We can use LoadStuffQA to create a chain that takes input documents and a question,
	// stuffs all the documents into the prompt of the llm and returns an answer to the
	// question. It is suitable for a small number of documents.
	stuffQAChain := chains.LoadStuffQA(llm)

		
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": searchResults,
		"question":        question,
	})
	if err != nil {
		return "",err
	}
	fmt.Println("RAG stuffed QA answer: ", answer)
	

		s,ok := answer["text"].(string)
		if ok {
			result = s
		}

		fmt.Println("====final answer  (STUFFED QA)====\n", result)
	
		return result, nil


}




// Test shows bad result, probably will return here with langraph, cause langchain is not flexible enough to control one-by-one chain
func RefinedQA_RAG(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error) {
	
		//base_url := os.Getenv("AI_BASEURL")
		base_url := ai_url

		// Create an embeddings client using the specified API and embedding model
		llm, err := openai.New(
			openai.WithBaseURL(base_url),
			openai.WithAPIVersion("v1"),
			openai.WithToken(api_token),
			openai.WithModel("tiger-gemma-9b-v1-i1"),
			openai.WithEmbeddingModel("text-embedding-ada-002"),
		)
		if err != nil {
			return "", err
		}
	
		//docs we found in this store
		searchResults, err := embeddings.SemanticSearch(question, numOfResults, store)
		if err != nil {
			return "", err
		}

	
	// Another option is to use the refine documents chain for question answering. This
	// chain iterates over the input documents one by one, updating an intermediate answer
	// with each iteration. It uses the previous version of the answer and the next document
	// as context. The downside of this type of chain is that it uses multiple llm calls that
	// cant be done in parallel.
	refineQAChain := chains.LoadRefineQA(llm)
	answer, err := chains.Call(context.Background(), refineQAChain, map[string]any{
		"input_documents": searchResults,
		"question":        "Where did Ankush go to collage?",
	})
	if err != nil {
		return "",err
	}
	fmt.Println("RAG Refined QA answer: ", answer)
	

		s,ok := answer["text"].(string)
		if ok {
			result = s
		}

		fmt.Println("====final answer  (REFINED QA)====\n", result)
	
		return result, nil



}