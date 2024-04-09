package handler

import (
	"encoding/json"
	"errors"
	"github.com/vavelour/chat/pkg/pagination"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/handler/mapper"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/internal/handler/response"
	"github.com/vavelour/chat/pkg/http_utils/baseresponse"
)

const (
	usersReceived    = "users received"
	messagesReceived = "messages received"
	messageSent      = "message sent"
	messagesNotFound = "no messages found"
)

var errFailedGetSender = errors.New("failed to get sender")

//go:generate mockgen -source=private_handler.go -destination=mocks/private_service_mock.go

type PrivateService interface {
	SendPrivateMessage(m entities.Message) error
	GetPrivateMessages(sender, recipient string, limit, offset int) ([]entities.Message, error)
	ViewUsers(user string) ([]string, error)
}

type PrivateHandler struct {
	service  PrivateService
	validate *validator.Validate
}

func NewPrivateHAndler(s PrivateService, v *validator.Validate) *PrivateHandler {
	return &PrivateHandler{service: s, validate: v}
}

func (h *PrivateHandler) PrivateRoutes(router *chi.Mux, middlewares ...func(next http.Handler) http.Handler) {
	router.Route("/v1/private", func(r chi.Router) {
		for _, mw := range middlewares {
			r.Use(mw)
		}
		r.Get("/users", h.ViewUserList)
		r.Get("/messages", h.ShowPrivateMessages)
		r.Post("/messages", h.SendPrivateMessage)
	})
}

// SendPrivateMessage @summary		Отправка приватного сообщения
//
//	@description	Отправляет приватное сообщение от имени отправителя указанному получателю.
//	@tags			private
//	@accept			json
//	@produce		json
//
//	@Security		BasicAuth
//
//	@param			username	query		string								true	"Имя получателя"
//	@param			requestBody	body		request.SendPrivateMessageRequest	true	"Данные сообщения"
//	@success		200			{object}	response.SendPrivateMessageResponse	"Сообщение успешно отправлено"
//	@failure		400			{object}	baseresponse.ResponseError			"Неверный запрос"
//	@router			/v1/private/messages [post]
func (h *PrivateHandler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	var input request.SendPrivateMessageRequest

	sender, ok := r.Context().Value("Sender").(string)
	if !ok {
		return
	}

	recipient := r.URL.Query().Get("username")

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	input.Sender = sender
	input.Recipient = recipient

	err := input.Validate(h.validate)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.service.SendPrivateMessage(mapper.SendPrivateMessageRequestToEntities(input))
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response.SendPrivateMessageResponse{Response: messageSent})
}

// ShowPrivateMessages @summary		Получение приватных сообщений
//
//	@description	Получает приватные сообщения между отправителем и получателем с заданным лимитом и смещением.
//	@tags			private
//	@accept			json
//	@produce		json
//
//	@Security		BasicAuth
//
//	@param			username	query		string								true	"Имя отправителя/получателя"
//	@param			requestBody	body		request.ShowPrivateMessageRequest	true	"Параметры запроса сообщений"
//	@success		200			{object}	response.ShowPrivateMessageResponse	"Сообщения успешно получены"
//	@failure		400			{object}	baseresponse.ResponseError			"Неверный запрос"
//	@failure		500			{object}	baseresponse.ResponseError			"Ошибка при получении сообщений"
//	@failure		416			{object}	baseresponse.ResponseError			"Запрос содержит невыполнимый диапазон"
//	@router			/v1/private/messages [get]
func (h *PrivateHandler) ShowPrivateMessages(w http.ResponseWriter, r *http.Request) {
	var input request.ShowPrivateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	sender, ok := r.Context().Value("Sender").(string)
	if !ok {
		baseresponse.ReturnErrorResponse(w, r, http.StatusInternalServerError, errFailedGetSender)
		return
	}

	recipient := r.URL.Query().Get("username")

	input.Sender = sender
	input.Recipient = recipient

	err := input.Validate(h.validate)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	messages, err := h.service.GetPrivateMessages(input.Sender, input.Recipient, input.Limit, input.Offset)
	if err != nil && !errors.Is(err, pagination.ErrOffsetRange) {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	} else if errors.Is(err, pagination.ErrOffsetRange) {
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, mapper.PrivateMessageEntitiesToResponse(messagesNotFound, messages))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, mapper.PrivateMessageEntitiesToResponse(messagesReceived, messages))
}

// ViewUserList @summary		Получение списка пользователей, от которых поступали сообщения
//
//	@description	Получает список пользователей.
//	@tags			private
//	@accept			json
//	@produce		json
//
//	@Security		BasicAuth
//
//	@success		200	{object}	response.ViewUserListResponse	"Список пользователей успешно получен"
//	@failure		400	{object}	baseresponse.ResponseError		"Неверный запрос"
//	@failure		500	{object}	baseresponse.ResponseError		"Ошибка при получении списка пользователей"
//	@router			/v1/private/users [get]
func (h *PrivateHandler) ViewUserList(w http.ResponseWriter, r *http.Request) {
	var input request.ViewUserListRequest

	username, ok := r.Context().Value("Sender").(string)
	if !ok {
		baseresponse.ReturnErrorResponse(w, r, http.StatusInternalServerError, errFailedGetSender)
		return
	}

	input.Username = username

	err := input.Validate(h.validate)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	users, err := h.service.ViewUsers(input.Username)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, mapper.UserListEntitiesToResponse(usersReceived, users))
}
