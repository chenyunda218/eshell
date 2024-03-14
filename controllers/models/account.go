package models

type CreateSessionRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type CreateDashboardRequest struct {
	Name string `json:"name"`
}

type Dashboard struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Visualization struct {
	Type  string `json:"type"`
	DSL   any    `json:"dsl"`
	Views any    `json:"views"`
}

type UpdateDashboardRequest struct {
	Name string `json:"name"`
}
