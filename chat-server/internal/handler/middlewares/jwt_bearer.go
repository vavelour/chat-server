package middlewares

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/pkg/http_utils/baseresponse"
	"net/http"
	"strings"
)

//go:generate mockgen -source=jwt_bearer.go -destination=mocks/jwt_identity_mock.go

type JWTBearerService interface {
	UserIdentity(token interface{}) (string, error)
}

type JWTUserIdentity struct {
	service  JWTBearerService
	validate *validator.Validate
}

func NewJWTUserIdentity(s JWTBearerService, v *validator.Validate) *JWTUserIdentity {
	return &JWTUserIdentity{service: s, validate: v}
}

func (h *JWTUserIdentity) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, errMissingAuth)
			return
		}

		headerParts := strings.SplitN(header, " ", credentialsNumber)
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, errInvalidAuth)
			return
		}

		req := request.JWTLogInRequest{Token: headerParts[1]}

		err := req.Validate(h.validate)
		if err != nil {
			baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, errEmptyLoginAndPassword)
			return
		}

		username, err := h.service.UserIdentity(req.Token)
		if err != nil {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), "Sender", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
