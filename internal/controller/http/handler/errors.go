package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrCodeUnauthorized        = errors.New("UNAUTHORIZED")
	ErrCodeNotFound            = errors.New("NOT_FOUND")
	ErrCodeBadRequest          = errors.New("BAD_REQUEST")
	ErrCodeForbidden           = errors.New("FORBIDDEN")
	ErrCodeConflict            = errors.New("CONFLICT")
	ErrCodeUnprocessableEntity = errors.New("UNPROCESSABLE_ENTITY")
	ErrCodeTooManyRequests     = errors.New("TOO_MANY_REQUESTS")
	ErrCodeInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
	ErrCodeBadGateway          = errors.New("BAD_GATEWAY")
	ErrCodeServiceUnavailable  = errors.New("SERVICE_UNAVAILABLE")
	ErrCodeGatewayTimeout      = errors.New("GATEWAY_TIMEOUT")
	ErrCodeCircuitBreakerOpen  = errors.New("CIRCUIT_BREAKER_OPEN")
)

type CustomWebError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CreateWebError(status int, code, message string) *CustomWebError {
	return &CustomWebError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func GlobalErrorHandler(err error) *CustomWebError {

	var webErr *CustomWebError

	switch {

	// 4xx Client Errors
	case errors.Is(err, ErrCodeNotFound):
		webErr = CreateWebError(http.StatusNotFound, ErrCodeNotFound.Error(), "Resource not found.")
	case errors.Is(err, ErrCodeBadRequest):
		webErr = CreateWebError(http.StatusBadRequest, ErrCodeBadRequest.Error(), "The request is invalid.")
	case errors.Is(err, ErrCodeUnauthorized):
		webErr = CreateWebError(http.StatusUnauthorized, ErrCodeUnauthorized.Error(), "Authentication is required.")
	case errors.Is(err, ErrCodeForbidden):
		webErr = CreateWebError(http.StatusForbidden, ErrCodeForbidden.Error(), "You are not allowed to perform this action.")
	case errors.Is(err, ErrCodeConflict):
		webErr = CreateWebError(http.StatusConflict, ErrCodeConflict.Error(), "The request conflicts with the current resource state.")
	case errors.Is(err, ErrCodeUnprocessableEntity):
		webErr = CreateWebError(http.StatusUnprocessableEntity, ErrCodeUnprocessableEntity.Error(), "The request could not be processed.")
	case errors.Is(err, ErrCodeTooManyRequests):
		webErr = CreateWebError(http.StatusTooManyRequests, ErrCodeTooManyRequests.Error(), "Too many requests.")

	// 5xx Server Errors
	case errors.Is(err, ErrCodeBadGateway):
		webErr = CreateWebError(http.StatusBadGateway, ErrCodeBadGateway.Error(), "A downstream service returned an invalid response.")
	case errors.Is(err, ErrCodeServiceUnavailable):
		webErr = CreateWebError(http.StatusServiceUnavailable, ErrCodeServiceUnavailable.Error(), "The service is temporarily unavailable.")
	case errors.Is(err, ErrCodeGatewayTimeout):
		webErr = CreateWebError(http.StatusGatewayTimeout, ErrCodeGatewayTimeout.Error(), "A downstream service timed out.")

	// Circuit Breaker Errors
	case errors.Is(err, ErrCodeCircuitBreakerOpen):
		webErr = CreateWebError(http.StatusServiceUnavailable, ErrCodeCircuitBreakerOpen.Error(), "The service is temporarily unavailable.")

	default:
		webErr = CreateWebError(http.StatusInternalServerError, ErrCodeInternalServerError.Error(), "An unexpected error occurred.")
	}

	return webErr
}

func WriteJSONError(c *gin.Context, webErr *CustomWebError) {
	c.JSON(webErr.Status, Response[any]{
		Success: false,
		Error: &APIError{
			Status:  webErr.Status,
			Code:    webErr.Code,
			Message: webErr.Message,
		},
	})
	//c.Abort()
}
