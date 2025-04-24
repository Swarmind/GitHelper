# agent

This package provides a conversational agent that can access a semantic search tool to retrieve information from a database or collection.

## File structure
- agent.go
- semantic_search_test.go
- superagent.go
- tools/tools.go

## Code summary
The package contains two main functions: `OneShotRun` and `RunNewAgent`.

### OneShotRun
This function takes a prompt, an LLM model, and an optional history state as input. It then uses the model to generate a response, which is returned as a string.

The function first initializes the agent's state and the initial state of the conversation. If a history state is provided, it appends the history state to the initial state.

Next, the function creates a message graph workflow, which includes nodes for the agent and the semantic search tool. The workflow is then compiled and invoked with the initial state.

Finally, the function extracts the last message from the response and returns its content as a string.

### RunNewAgent
This function creates a new chat session graph, initializes an LLM model, and runs the `CreateThread` function with the provided prompt, model, and collection names.

The function first checks if the base URL is empty. If it is, it creates a new LLM using the provided AI token and model. Otherwise, it creates a new LLM using the provided AI token, model, and base URL.

Next, the function calls the `CreateThread` function to create a new thread with the provided prompt, model, and collection names. The `CreateThread` function returns the updated history state and the call result.

Finally, the function returns the chat session graph, output text, and any errors.

### Other functions
The package also contains several other functions, such as `OnePunch`, `RunThread`, `CreateThread`, `CreateMessageContentAi`, `createMessageContentSystem`, `CreateMessageContentHuman`, `CreateGenericLLM`, and `ContinueAgent`. These functions are used to support the main functions and provide additional functionality.

### Tools
The package also includes a `tools` subpackage, which contains a `SemanticSearchTool` struct. This struct is used to execute semantic search queries.

### Configuration
The package does not appear to have any configuration files or environment variables.

### Edge cases
The package does not appear to have any edge cases or special handling for specific situations.