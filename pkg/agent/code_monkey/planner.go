package code_monekey

import (
	"context"

	graph "github.com/JackBekket/langgraphgo/graph/stategraph"
)

type Plan struct {
    Task string
    Steps []string
    PlanString string
}

type ReWOO struct {
    Task string
    PlanString string
    Steps []string
    Results map[string]string
    Result string
}

// get_plan - функция для работы со State
func get_plan(ctx context.Context, state graph.State) graph.State {
    task := state.Get("task").(string)
    result := planner.invoke(map[string]string{"task": task})
    matches := re.findall(regex_pattern, result.content)
    return ReWOO{
        "steps": matches,
        "plan_string": result.content,
    }
}