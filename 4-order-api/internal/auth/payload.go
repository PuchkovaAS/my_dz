package auth

type GetAuthCodeRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type GetAuthCodeResponse struct {
	SessionId string `json:"sessionID"`
}

type GetTokenRequest struct {
	SessionId string `json:"sessionID" validate:"required"`
	Code      uint   `json:"code"      validate:"required,gt=0"`
}

type GetTokenResponse struct {
	Token string `json:"token"`
}
