package rag

import (
	"context"
	"fmt"

	"github.com/JackBekket/hellper/lib/embeddings"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores"
)

// Hardcoded func which call RAG with doc & code metadata
func RagReflexia(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error) {
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

		// First step -- get the docs
		filters := map[string]any{
			"type": "doc",
		}
		
		option := vectorstores.WithFilters(filters)  
		Documentation, err := embeddings.SemanticSearch(question, numOfResults, store, option)
		if err != nil {
			return "", err
		}
		//----------------------  we got documentation 
		// step 2 -- get the code and add to documentation
				// filteres
				filters = map[string]any{
					"type": "code",
				}
				
				option = vectorstores.WithFilters(filters)  
				CodeContent, err := embeddings.SemanticSearch(question, numOfResults, store, option)
				if err != nil {
					return "", err
				}

			InputDocs := append(Documentation, CodeContent...)    // 
	
	stuffQAChain := chains.LoadStuffQA(llm)

		
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": InputDocs,
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
		
		fmt.Println("(RAG REFLEXIA DOC & CODE)====final answer====\n", result)
		return result, nil
}


func RagWithOptions(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore,option ...vectorstores.Option)  (result string, err error) {
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
		result, err = chains.Run(
			context.Background(),
			chains.NewRetrievalQAFromLLM(	// stuffed QA
				llm,
				vectorstores.ToRetriever(store, numOfResults,option...),
			),
			question,
			chains.WithMaxTokens(8192),
		)
		if err != nil {
			return "", err
		}
	fmt.Println("(RAG WITH ONLY FILTERES AS OPTION)====final answer====\n", result)
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




func RefinedQA_RAG(question string, ai_url string, api_token string, numOfResults int, store vectorstores.VectorStore) (result string, err error) {
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
		"question":        question,
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