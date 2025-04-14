package code_monekey

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

