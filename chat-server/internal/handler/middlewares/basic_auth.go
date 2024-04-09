package middlewares

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/vavelour/chat/internal/handler/mapper"

	"github.com/vavelour/chat/internal/handler/request"

	"github.com/vavelour/chat/pkg/http_utils/baseresponse"
)

const (
	credentialsNumber = 2
)

var (
	errMissingAuth           = errors.New("authorization header is missing")
	errInvalidAuth           = errors.New("invalid authorization header")
	errEmptyLoginAndPassword = errors.New("login and password cannot be empty")
)

//go:generate mockgen -source=basic_auth.go -destination=mocks/basic_identity_mock.go

type BasicAuthService interface {
	UserIdentity(usr interface{}) (string, error)
}

type BasicUserIdentity struct {
	service  BasicAuthService
	validate *validator.Validate
}

func NewBasicUserIdentity(s BasicAuthService, v *validator.Validate) *BasicUserIdentity {
	return &BasicUserIdentity{service: s, validate: v}
}

func (h *BasicUserIdentity) Identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, errMissingAuth)
			return
		}

		headerParts := strings.SplitN(header, " ", credentialsNumber)
		if len(headerParts) != 2 || headerParts[0] != "Basic" || headerParts[1] == "" {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, errInvalidAuth)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(headerParts[1])
		if err != nil {
			baseresponse.ReturnErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", credentialsNumber)
		if len(credentials) != credentialsNumber {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		req := request.BasicAuthLogInRequest{Username: credentials[0], Password: credentials[1]}

		err = req.Validate(h.validate)
		if err != nil {
			baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, errEmptyLoginAndPassword)
			return
		}

		username, err := h.service.UserIdentity(mapper.BasicLogInRequestToEntities(req))
		if err != nil {
			baseresponse.ReturnErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), "Sender", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
