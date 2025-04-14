package code_monekey

import "fmt"

// ReWOO represents a dictionary-like data structure for storing state information.
//type ReWOO map[string]interface{}

// _get_current_task determines the index of the next task to be executed.
func _get_current_task(state ReWOO) int {
    if _, ok := state["results"]; !ok || state["results"] == nil {
        return 1
    }
    if len(state["results"].(map[string]interface{})) == len(state["steps"].([]interface{})) {
        return -1
    }
    return len(state["results"].(map[string]interface{})) + 1
}

// tool_execution executes the tools of a given plan.
func tool_execution(state ReWOO) (ReWOO, error) {
    _step := _get_current_task(state)
    if _step == -1 {
        return nil, fmt.Errorf("no more tasks to execute")
    }
    _, step_name, tool, tool_input := state["steps"].([]interface{})[_step-1]
    _results := state["results"].(map[string]interface{})
    for k, v := range _results {
        tool_input = tool_input.(string)
        tool_input = string(v.(string))
    }
    var result string
    switch tool.(string) {
    case "Google":
        // Implement search.invoke function
        result = searchInvoke(tool_input)
    case "LLM":
        // Implement model.invoke function
        result = modelInvoke(tool_input)
    default:
        return nil, fmt.Errorf("unsupported tool")
    }
    _results[step_name] = result
    return ReWOO{"results": _results}, nil
}