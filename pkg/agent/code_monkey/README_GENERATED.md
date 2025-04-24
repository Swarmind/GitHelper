## Package: code_monkey

Imports:
- context
- fmt
- strings
- github.com/JackBekket/GitHelper/pkg/agent/rag
- github.com/rs/zerolog/log
- github.com/tmc/langchaingo/llms

External data, input sources:
- The code uses an LLM (Large Language Model) to process tasks. The LLM is initialized and used through the `llms` package.
- The code interacts with a toolchain to execute steps in a workflow. The toolchain is initialized and used through the `InitializeChain()` function.

TODOs:
- There are no TODO comments in the provided code.

### ToolExecution function
This function is responsible for executing a tool or LLM based on the current task in a workflow. It takes the current state of the workflow and the context as input.

First, it replaces any previously executed tool names in the tool input with their corresponding results. Then, it checks if the current tool is an LLM or a regular tool. If it's an LLM, it generates content using the LLM's `GenerateContent` function. If it's a regular tool, it executes the tool using the toolchain's `Invoke` function.

Finally, it updates the state of the workflow with the result of the tool or LLM execution and returns the updated state.

### Graph.go
- The code defines a struct called LLMContext, which holds an OpenAI LLM and a list of tools.
- It implements a method called OneShotRun, which takes a context and a prompt as input.
- The method creates a state graph with nodes for planning, tool execution, and solving.
- It then compiles the graph and invokes it with the given prompt.
- The result of the invocation is returned as a string.

This code appears to be part of a system that uses an LLM to generate plans and execute tools to solve tasks. The state graph is used to manage the workflow of the system, and the OneShotRun method provides a way to run the system with a given prompt.

### Planner.go
- The code defines a ReWOO struct, which represents a task, its plan, steps, results, and a final result.
- It also defines a ReWOOStep struct, which represents a single step in the plan.
- The code includes functions to generate plans, extract steps from the plan, and manage the state of the task.

The GetPlan function takes a task as input and uses an LLM to generate a plan. The plan is then parsed to extract steps, which are stored in the ReWOO struct. The function also handles the state of the task, keeping track of the current step and the results of previous steps.

The getToolDesc function takes a list of tools and returns a string containing the description of each tool. This function is used to provide the LLM with information about the available tools.

The getCurrentTask function determines the current task being worked on based on the number of results obtained so far.

### Solver.go
- The Solve function takes a context, a state of type ReWOO, and returns the updated state and an error. It constructs a prompt for the LLM by combining the plan and the task from the input state. Then, it uses the LLM to generate a response and updates the state with the response.
- The Route function takes a context and a state of type ReWOO and returns a string. It checks the current task index in the state and returns GraphSolveName if it's -1, otherwise it returns GraphToolName.

### Toolchain.go
- The code defines a state graph workflow for generating tool calls and performing semantic search.
- The workflow consists of two nodes: generate_call and semanticSearch.
- The generate_call node generates a tool call based on the input prompt.
- The semanticSearch node performs a similarity search in a vectorstore to find relevant documents.
- The workflow starts with the generate_call node and then conditionally moves to the semanticSearch node based on the tool call.
- The code also includes functions for initializing the language model and executing the workflow.

- The initializeModel function initializes the language model using the API_TOKEN, AI_URL, and MODEL environment variables.
- The InitializeChain function creates the state graph workflow and compiles it into a runnable application.
- The semanticSearch function performs the similarity search in the vectorstore.
- The generateCall function generates the tool call based on the input prompt.
- The whichTool function determines which tool node to use based on the tool call.

pkg/agent/code_monkey/executor.go
pkg/agent/code_monkey/graph.go
pkg/agent/code_monkey/planner.go
pkg/agent/code_monkey/solver.go
pkg/agent/code_monkey/toolchain.go
pkg/agent/code_monkey/planner_test.go
