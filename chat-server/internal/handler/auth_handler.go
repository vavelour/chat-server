package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vavelour/chat/internal/handler/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/pkg/http_utils/baseresponse"
)

const (
	userCreated = "user created"
)

//go:generate mockgen -source=auth_handler.go -destination=mocks/auth_service_mock.go

type AuthService interface {
	CreateUser(username, password string) (string, error)
}

type AuthHandler struct {
	service  AuthService
	validate *validator.Validate
}

func NewAuthHandler(s AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{service: s, validate: v}
}

func (h *AuthHandler) AuthRoutes(router *chi.Mux, middlewares ...func(next http.Handler) http.Handler) {
	router.Route("/v1/auth", func(r chi.Router) {
		for _, mw := range middlewares {
			r.Use(mw)
		}
		r.Post("/register", h.Register)
	})
}

// Register @summary		Регистрация нового пользователя
//
//	@description	Регистрирует нового пользователя с заданным именем пользователя и паролем.
//	@tags			auth
//	@accept			json
//	@produce		json
//	@param			requestBody	body		request.RegisterRequest		true	"Данные нового пользователя"
//	@success		201			{object}	response.RegisterResponse	"Пользователь успешно создан"
//	@failure		400			{object}	baseresponse.ResponseError	"Неверный запрос"
//	@failure		409			{object}	baseresponse.ResponseError	"Пользователь уже существует"
//	@router			/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	err := input.Validate(h.validate)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.CreateUser(input.Username, input.Password)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusConflict, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response.RegisterResponse{Response: userCreated, Token: token})
}
