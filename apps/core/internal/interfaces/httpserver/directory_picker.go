package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/directorypicker"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/httpserver/api"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

const maxDirectoryPickerRequestBody = 1 << 20

var errInvalidDirectoryPickerRequest = errors.New("invalid directory picker request")

func (server *Server) selectDirectoryHandler(c *echo.Context) error {
	request, err := decodeSelectDirectoryRequest(c)
	if err != nil {
		return server.writeDirectoryPickerError(c, err)
	}
	if request.Purpose != api.ProjectLocation {
		return server.writeDirectoryPickerError(
			c,
			fmt.Errorf("%w: purpose must be project_location", errInvalidDirectoryPickerRequest),
		)
	}

	path, selected, err := server.directoryPicker.Select(c.Request().Context())
	if err != nil {
		return server.writeDirectoryPickerError(c, err)
	}
	response := api.SelectDirectoryResponse{Selected: selected}
	if selected {
		response.Path = &path
	}
	return c.JSON(http.StatusOK, response)
}

func decodeSelectDirectoryRequest(c *echo.Context) (api.SelectDirectoryRequest, error) {
	request := c.Request()
	request.Body = http.MaxBytesReader(c.Response(), request.Body, maxDirectoryPickerRequestBody)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var body api.SelectDirectoryRequest
	if err := decoder.Decode(&body); err != nil {
		return api.SelectDirectoryRequest{}, fmt.Errorf(
			"%w: request must contain one valid JSON object",
			errInvalidDirectoryPickerRequest,
		)
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return api.SelectDirectoryRequest{}, fmt.Errorf(
			"%w: request must not contain trailing JSON",
			errInvalidDirectoryPickerRequest,
		)
	}
	return body, nil
}

func (server *Server) writeDirectoryPickerError(c *echo.Context, err error) error {
	status := http.StatusInternalServerError
	code := api.ErrorCodeInternalError
	message := "Core could not complete the request"
	recoverable := false

	switch {
	case errors.Is(err, errInvalidDirectoryPickerRequest):
		status = http.StatusBadRequest
		code = api.ErrorCodeValidationFailed
		message = err.Error()
		recoverable = true
	case errors.Is(err, directorypicker.ErrUnavailable):
		status = http.StatusServiceUnavailable
		code = api.ErrorCodeDirectoryPickerUnavailable
		message = "Native directory picker is unavailable; enter the absolute path manually"
		recoverable = true
	default:
		server.logger.Error("directory picker request failed", zap.Error(err))
	}

	return c.JSON(status, api.ErrorResponse{Error: api.Error{
		Code:        code,
		Message:     message,
		Recoverable: recoverable,
	}})
}
