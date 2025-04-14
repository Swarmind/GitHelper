package code_monekey

import (
	"fmt"
	"regexp"
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

func getPlan(task string) Plan {
    planString := fmt.Sprintf("Plan: %s\n", task)
    regexPattern := regexp.MustCompile(`Plan:\s*(.+)\s*(#E\d+)\s*=\s*(\w+)\s*\[([^\]]+)\]`)
    planString = planString + "#E1 = step1 [some content]\n"
    planString = planString + "#E2 = step2 [some content]\n"
    planString = planString + "#E3 = step3 [some content]\n"
    planString = planString + "#E4 = step4 [some content]\n"
    planString = planString + "#E5 = step5 [some content]\n"
    matches := regexPattern.FindAllStringSubmatch(planString, -1)
    steps := make([]string, len(matches))
    for i, match := range matches {
        steps[i] = match[2]
    }
    return Plan{Task: task, Steps: steps, PlanString: planString}
}