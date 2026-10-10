package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/go-basic-clean-template/internal/controller/http/handler"
	coreHealth "github.com/pinedadaniel/go-basic-clean-template/internal/core/usecase/health"
)

type healthUseCaseStub struct{}

func (healthUseCaseStub) Execute(context.Context) (coreHealth.Output, error) {
	return coreHealth.Output{}, nil
}

func TestGetMissingCallerIDReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", New(healthUseCaseStub{}).Get)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}

	var response handler.Response[any]
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if response.Error == nil {
		t.Fatal("response error is nil")
	}
	if response.Error.Code != ErrCodeCallerIDRequired {
		t.Errorf("error code = %q, want %q", response.Error.Code, ErrCodeCallerIDRequired)
	}
	if response.Error.Message != "Header X-Caller-Id is required." {
		t.Errorf("error message = %q, want safe caller ID validation message", response.Error.Message)
	}
}
