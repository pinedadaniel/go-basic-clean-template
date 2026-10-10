package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGlobalErrorHandlerUsesSafeInternalErrorResponse(t *testing.T) {
	for _, testErr := range []error{
		errors.New("database password leaked"),
		ErrCodeInternalServerError,
	} {
		webErr := GlobalErrorHandler(testErr)

		if webErr.Status != http.StatusInternalServerError {
			t.Errorf("status for %v = %d, want %d", testErr, webErr.Status, http.StatusInternalServerError)
		}
		if webErr.Code != ErrCodeInternalServerError.Error() {
			t.Errorf("code for %v = %q, want %q", testErr, webErr.Code, ErrCodeInternalServerError.Error())
		}
		if webErr.Message != "An unexpected error occurred." {
			t.Errorf("message for %v = %q, want a generic public message", testErr, webErr.Message)
		}
	}
}

func TestGlobalErrorHandlerDoesNotEchoRetryAfterRequestHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("Retry-After", "secret-client-value")

	webErr := GlobalErrorHandler(ErrCodeTooManyRequests)
	WriteJSONError(ctx, webErr)

	if got := recorder.Header().Get("Retry-After"); got != "" {
		t.Errorf("Retry-After response header = %q, want omitted", got)
	}
	if recorder.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
}
