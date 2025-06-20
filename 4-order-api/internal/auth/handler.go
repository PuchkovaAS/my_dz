package auth

import (
	"4-order-api/configs"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/request"
	"4-order-api/pkg/response"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
	*jwt.JWT
}

type AuthHandler struct {
	*configs.Config
	*AuthService
	*jwt.JWT
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		JWT:         deps.JWT,
	}
	router.HandleFunc("POST /auth/get_code", handler.GetCode())
	router.HandleFunc("POST /auth/get_token", handler.GetToken())
}

func (handler *AuthHandler) GetCode() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[GetAuthCodeRequest](&w, req)
		if err != nil {
			return
		}

		user, err := handler.AuthService.CreateSession(body.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		data := GetAuthCodeResponse{
			SessionId: user.SessionId,
		}
		response.Json(w, data, http.StatusCreated)
	}
}

func (handler *AuthHandler) GetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[GetTokenRequest](&w, req)
		if err != nil {
			return
		}

		userPhone, err := handler.AuthService.VerifyCode(
			body.SessionId,
			uint(body.Code),
		)
		if err != nil || userPhone != "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tokenString, err := handler.JWT.Create(
			jwt.JWTData{Phone: userPhone},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := GetTokenResponse{
			Token: tokenString,
		}
		response.Json(w, data, http.StatusCreated)
	}
}
