package handler

import (
	"encoding/json"
	"errors"
	"github.com/vavelour/chat/pkg/pagination"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/vavelour/chat/internal/domain/entities"

	"github.com/go-chi/chi/v5"
	"github.com/vavelour/chat/internal/handler/mapper"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/internal/handler/response"
	"github.com/vavelour/chat/pkg/http_utils/baseresponse"
)

//go:generate mockgen -source=public_handler.go -destination=mocks/public_service_mock.go

type PublicService interface {
	SendPublicMessage(m entities.Message) error
	GetPublicMessages(limit, offset int) ([]entities.Message, error)
}

type PublicHandler struct {
	service  PublicService
	validate *validator.Validate
}

func NewPublicHandler(h PublicService, v *validator.Validate) *PublicHandler {
	return &PublicHandler{service: h, validate: v}
}

func (h *PublicHandler) PublicRoutes(router *chi.Mux, middlewares ...func(next http.Handler) http.Handler) {
	router.Route("/v1/public", func(r chi.Router) {
		for _, mw := range middlewares {
			r.Use(mw)
		}
		r.Get("/messages", h.ShowPublicMessages)
		r.Post("/messages", h.SendPublicMessage)
	})
}

// SendPublicMessage @summary		Отправка сообщения в публичный чат
//
//	@description	Отправляет сообщение в публичный чат от имени пользователя.
//	@tags			public
//	@accept			json
//	@produce		json
//
//	@Security		BasicAuth
//
//	@param			requestBody	body		request.SendPublicMessageRequest	true	"Данные сообщения"
//	@success		200			{object}	response.SendPublicMessageResponse	"Сообщение успешно отправлено"
//	@failure		500			{object}	baseresponse.ResponseError			"Ошибка при отправке сообщения"
//	@failure		400			{object}	baseresponse.ResponseError			"Неверный запрос"
//	@router			/v1/public/messages [post]
func (h *PublicHandler) SendPublicMessage(w http.ResponseWriter, r *http.Request) {
	var input request.SendPublicMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	sender, ok := r.Context().Value("Sender").(string)
	if !ok {
		baseresponse.ReturnErrorResponse(w, r, http.StatusInternalServerError, errFailedGetSender)
		return
	}

	input.Sender = sender

	err := input.Validate(h.validate)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.service.SendPublicMessage(mapper.SendPublicMessageRequestToEntities(input))
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response.SendPublicMessageResponse{Response: messageSent})
}

// ShowPublicMessages @summary		Получение сообщений из публичного чата
//
//	@description	Получает сообщения из публичного чата с заданным лимитом и смещением.
//	@tags			public
//	@accept			json
//	@produce		json
//
//	@Security		BasicAuth
//
//	@param			requestBody	body		request.ShowPublicMessageRequest	true	"Параметры запроса сообщений"
//	@success		200			{object}	response.ShowPublicMessageResponse	"Сообщения успешно получены"
//	@failure		400			{object}	baseresponse.ResponseError			"Неверный запрос"
//	@failure		416			{object}	baseresponse.ResponseError			"Запрос содержит невыполнимый диапазон"
//	@router			/v1/public/messages [get]
func (h *PublicHandler) ShowPublicMessages(w http.ResponseWriter, r *http.Request) {
	var input request.ShowPublicMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	err := input.Validate(h.validate)
	if err != nil {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	messages, err := h.service.GetPublicMessages(input.Limit, input.Offset)
	if err != nil && !errors.Is(err, pagination.ErrOffsetRange) {
		baseresponse.ReturnErrorResponse(w, r, http.StatusBadRequest, err)
		return
	} else if errors.Is(err, pagination.ErrOffsetRange) {
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, mapper.PublicMessageEntitiesToResponse(messagesNotFound, messages))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, mapper.PublicMessageEntitiesToResponse(messagesReceived, messages))
}
