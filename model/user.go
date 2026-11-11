package model

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}
