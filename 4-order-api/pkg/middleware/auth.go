package middleware

import (
	"4-order-api/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextPhoneKey key = "ContextPhoneKey"
)

func writeUnauthted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, jwtService jwt.JWT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerAuth := r.Header.Get("Authorization")
		if !strings.HasPrefix(headerAuth, "Bearer ") {
			writeUnauthted(w)
			return
		}
		token := strings.TrimPrefix(headerAuth, "Bearer ")
		isValid, data := jwtService.Parse(token)
		if !isValid {
			writeUnauthted(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		reqst := r.WithContext(ctx)
		next.ServeHTTP(w, reqst)
	})
}
