package verify

import (
	"3-validation-api/configs"
	"3-validation-api/internal/storage"
	"3-validation-api/pkg/hashing"
	"3-validation-api/pkg/linker"
	"3-validation-api/pkg/mailer"
	"3-validation-api/pkg/request"
	"3-validation-api/pkg/response"
	"net/http"
	"strings"
)

type VerifyHandlerDeps struct {
	*configs.Config
}

type VerifyHandler struct {
	*configs.Config
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[SendRequest](&w, req)
		if err != nil {
			return
		}

		hashString := hashing.GetHashString(body.Email)
		urlLink := linker.GetHashUrl("http://localhost:8081/verify", hashString)

		msg := mailer.SendMsg{
			AddressTo:  body.Email,
			SubjectMsg: "Verification",
			TextMsg:    urlLink,
		}
		err = mailer.SendEmail(&handler.Config.Email, &msg)
		if err != nil {

			data := SendResponse{
				Status: "email did`t send",
			}
			response.Json(w, data, http.StatusAccepted)
			return

		}

		storage.AddEmailHash(
			handler.Config.Storage.Path,
			body.Email,
			hashString,
		)

		data := SendResponse{
			Status: "Ok",
		}

		response.Json(w, data, http.StatusAccepted)
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hashString := strings.TrimPrefix(req.URL.Path, "/verify/")

		var data VerifyResponse
		if isVerify := storage.CheckHash(
			handler.Config.Storage.Path,
			hashString); isVerify {
			data = VerifyResponse{
				Status: true,
			}
		} else {
			data = VerifyResponse{
				Status: false,
			}
		}

		response.Json(w, data, http.StatusAccepted)
	}
}
