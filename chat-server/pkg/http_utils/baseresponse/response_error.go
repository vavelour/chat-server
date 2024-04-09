package baseresponse

import (
	"net/http"

	"github.com/go-chi/render"
)

type ResponseError struct {
	Err string `json:"error"`
}

func ReturnErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, ResponseError{Err: err.Error()})
}
