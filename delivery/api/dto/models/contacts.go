package models

type SendContactsRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Message string `json:"message"`
}
