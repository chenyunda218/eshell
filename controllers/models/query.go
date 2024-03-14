package models

type Search struct {
	Size int `json:"size"`
}

type Query struct {
}

type CreateAccountRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
