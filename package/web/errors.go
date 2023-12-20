package web

import (
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
)

type responseError struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time_error"`
}

func BadRequestResponse(res *restful.Response, message string) {
	_ = res.WriteHeaderAndJson(http.StatusBadRequest, responseError{
		Message: message,
		Time:    time.Now(),
	}, restful.MIME_JSON)
}
