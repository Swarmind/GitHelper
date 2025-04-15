package code_monkey_test

import (
	"os"
	"testing"

	"github.com/JackBekket/GitHelper/pkg/agent/rag/tools"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"

	codeMonkey "github.com/JackBekket/GitHelper/pkg/agent/code_monkey"
)

func TestPlanner(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	llm, err := openai.New(
		openai.WithToken(os.Getenv("API_TOKEN")),
		openai.WithModel(os.Getenv("MODEL")),
		openai.WithBaseURL(os.Getenv("AI_URL")),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
		t.Fatal(err)
	}
	tools, err := tools.GetTools()
	if err != nil {
		t.Fatal(err)
	}

	lc := codeMonkey.LLMContext{
		LLM:   llm,
		Tools: &tools,
	}

	stateInterface, err := lc.GetPlan(t.Context(), codeMonkey.ReWOO{
		Task: `There is a specific code at 'Hellper' collection at vector store, related to a telegram bot api initialization.
I need to find a web documentation about that library and initialization sequence`,
	})
	if err != nil {
		t.Fatal(err)
	}
	state := stateInterface.(codeMonkey.ReWOO)

	t.Logf("Returned state: %+v", state)

	if len(state.Steps) == 0 {
		t.Fatal("empty steps")
	}
	if state.PlanString == "" {
		t.Fatal("empty plan string")
	}
}
