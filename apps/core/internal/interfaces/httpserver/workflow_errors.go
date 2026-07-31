package httpserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/httpserver/api"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

const maxWorkflowRequestBody = 1 << 20

func decodeWorkflowRequest[T any](c *echo.Context) (T, error) {
	var zero T
	request := c.Request()
	request.Body = http.MaxBytesReader(c.Response(), request.Body, maxWorkflowRequestBody)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	var body T
	if err := decoder.Decode(&body); err != nil {
		return zero, releasedomain.NewValidationError("request", "must contain one valid JSON object")
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return zero, releasedomain.NewValidationError("request", "must not contain trailing JSON")
	}
	return body, nil
}

func (server *Server) writeWorkflowError(c *echo.Context, err error) error {
	status := http.StatusInternalServerError
	code := api.ErrorCodeInternalError
	message := "Core could not complete the request"
	recoverable := false

	switch {
	case errors.Is(err, releasedomain.ErrValidation), errors.Is(err, checklistdomain.ErrValidation), errors.Is(err, projectdomain.ErrValidation):
		status, code, message, recoverable = http.StatusBadRequest, api.ErrorCodeValidationFailed, err.Error(), true
	case errors.Is(err, projectdomain.ErrNotFound):
		status, code, message, recoverable = http.StatusNotFound, api.ErrorCodeProjectNotFound, "Project not found", true
	case errors.Is(err, templatedomain.ErrNotFound):
		status, code, message, recoverable = http.StatusNotFound, api.ErrorCodeTemplateNotFound, "Release template not found", true
	case errors.Is(err, releasedomain.ErrNotFound):
		status, code, message, recoverable = http.StatusNotFound, api.ErrorCodeReleaseNotFound, "Release Goal not found", true
	case errors.Is(err, checklistdomain.ErrNotFound):
		status, code, message, recoverable = http.StatusNotFound, api.ErrorCodeChecklistItemNotFound, "Checklist item not found", true
	case errors.Is(err, releasedomain.ErrConflict):
		status, code, message, recoverable = http.StatusConflict, api.ErrorCodeReleaseGoalConflict, "A Release Goal of this type already exists", true
	case errors.Is(err, releasedomain.ErrNotReady):
		status, code, message, recoverable = http.StatusConflict, api.ErrorCodeReleaseNotReady, err.Error(), true
	case errors.Is(err, releasedomain.ErrInvalidTransition), errors.Is(err, checklistdomain.ErrInvalidTransition):
		status, code, message, recoverable = http.StatusConflict, api.ErrorCodeInvalidTransition, err.Error(), true
	case errors.Is(err, releasedomain.ErrStateConflict), errors.Is(err, checklistdomain.ErrReleaseState):
		status, code, message, recoverable = http.StatusConflict, api.ErrorCodeReleaseStateConflict, err.Error(), true
	default:
		server.logger.Error("Release workflow request failed", zap.Error(err))
	}

	return c.JSON(status, api.ErrorResponse{Error: api.Error{Code: code, Message: message, Recoverable: recoverable}})
}
