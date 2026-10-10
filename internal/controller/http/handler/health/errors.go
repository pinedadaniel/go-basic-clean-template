package health

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/go-basic-clean-template/internal/controller/http/handler"
)

const (
	ErrCodeCallerIDRequired      = "CALLER-ID-REQUIRED"
	ErrCodeDecodeJSONRequestBody = "DECODE-JSON-REQUEST-BODY"
	ErrCodeInvalidRequestHeaders = "INVALID-REQUEST-HEADERS"
)

var (
	ErrCallerIDRequired      = errors.New("header x-caller-id is required")
	ErrDecodeJSONRequestBody = errors.New("error decoding request body")
	ErrInvalidRequestHeaders = errors.New("request headers are invalid")
)

func handleError(ctx *gin.Context, err error) {
	var webErr *handler.CustomWebError

	switch {
	case errors.Is(err, ErrCallerIDRequired):
		webErr = handler.CreateWebError(http.StatusBadRequest, ErrCodeCallerIDRequired, "Header X-Caller-Id is required.")
	case errors.Is(err, ErrDecodeJSONRequestBody):
		webErr = handler.CreateWebError(http.StatusBadRequest, ErrCodeDecodeJSONRequestBody, "The request body is invalid.")
	case errors.Is(err, ErrInvalidRequestHeaders):
		webErr = handler.CreateWebError(http.StatusBadRequest, ErrCodeInvalidRequestHeaders, "The request headers are invalid.")
	default:
		webErr = handler.GlobalErrorHandler(err)
	}

	handler.WriteJSONError(ctx, webErr)
}
