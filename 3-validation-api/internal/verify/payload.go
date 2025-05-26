package verify

type SendRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SendResponse struct {
	Status string `json:"status"`
}

type VerifyResponse struct {
	Status bool `json:"status"`
}
