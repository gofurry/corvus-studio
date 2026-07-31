package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	projectapp "github.com/gofurry/corvus-studio/apps/core/internal/application/project"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/httpserver/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

const maxProjectRequestBody = 1 << 20

type ProjectService interface {
	Create(context.Context, projectapp.CreateCommand) (projectdomain.Project, error)
	Get(context.Context, projectdomain.ID) (projectdomain.Project, error)
	List(context.Context) ([]projectdomain.Project, error)
}

func (server *Server) createProjectHandler(c *echo.Context) error {
	request, err := decodeCreateProjectRequest(c)
	if err != nil {
		return server.writeProjectError(c, err)
	}
	description := ""
	if request.Description != nil {
		description = *request.Description
	}
	entity, err := server.projects.Create(c.Request().Context(), projectapp.CreateCommand{
		Name:        request.Name,
		Description: description,
		Location:    request.Location,
		SteamAppID:  request.SteamAppId,
		Language:    request.Language,
		Stage:       projectdomain.Stage(request.Stage),
	})
	if err != nil {
		return server.writeProjectError(c, err)
	}
	response, err := projectResponse(entity)
	if err != nil {
		return server.writeProjectError(c, err)
	}
	c.Response().Header().Set("Location", "/api/v1/projects/"+entity.ID.String())
	return c.JSON(http.StatusCreated, response)
}

func (server *Server) listProjectsHandler(c *echo.Context) error {
	entities, err := server.projects.List(c.Request().Context())
	if err != nil {
		return server.writeProjectError(c, err)
	}
	items := make([]api.Project, 0, len(entities))
	for _, entity := range entities {
		item, err := projectResponse(entity)
		if err != nil {
			return server.writeProjectError(c, err)
		}
		items = append(items, item)
	}
	return c.JSON(http.StatusOK, api.ProjectList{Items: items})
}

func (server *Server) getProjectHandler(c *echo.Context) error {
	id, err := projectdomain.ParseID(c.Param("project_id"))
	if err != nil {
		return server.writeProjectError(c, err)
	}
	entity, err := server.projects.Get(c.Request().Context(), id)
	if err != nil {
		return server.writeProjectError(c, err)
	}
	response, err := projectResponse(entity)
	if err != nil {
		return server.writeProjectError(c, err)
	}
	return c.JSON(http.StatusOK, response)
}

func decodeCreateProjectRequest(c *echo.Context) (api.CreateProjectRequest, error) {
	request := c.Request()
	request.Body = http.MaxBytesReader(c.Response(), request.Body, maxProjectRequestBody)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var body api.CreateProjectRequest
	if err := decoder.Decode(&body); err != nil {
		return api.CreateProjectRequest{}, projectdomain.NewValidationError(
			"request",
			"must contain one valid JSON object",
		)
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return api.CreateProjectRequest{}, projectdomain.NewValidationError(
			"request",
			"must not contain trailing JSON",
		)
	}
	return body, nil
}

func projectResponse(entity projectdomain.Project) (api.Project, error) {
	id, err := uuid.Parse(entity.ID.String())
	if err != nil || id.Version() != 7 {
		return api.Project{}, fmt.Errorf("encode Project ID %q: expected UUIDv7", entity.ID)
	}
	return api.Project{
		CreatedAt:   entity.CreatedAt.UTC(),
		Description: entity.Description,
		Id:          id,
		Language:    entity.Language,
		Location:    entity.Location,
		Name:        entity.Name,
		Stage:       api.ProjectStage(entity.Stage),
		Status:      api.ProjectStatus(entity.Status),
		SteamAppId:  cloneInt64(entity.SteamAppID),
		UpdatedAt:   entity.UpdatedAt.UTC(),
	}, nil
}

func (server *Server) writeProjectError(c *echo.Context, err error) error {
	status := http.StatusInternalServerError
	code := api.ErrorCodeInternalError
	message := "Core could not complete the request"
	recoverable := false

	switch {
	case errors.Is(err, projectdomain.ErrValidation):
		status = http.StatusBadRequest
		code = api.ErrorCodeValidationFailed
		message = err.Error()
		recoverable = true
	case errors.Is(err, projectdomain.ErrNotFound):
		status = http.StatusNotFound
		code = api.ErrorCodeProjectNotFound
		message = "Project not found"
		recoverable = true
	case errors.Is(err, projectdomain.ErrLocationConflict):
		status = http.StatusConflict
		code = api.ErrorCodeProjectLocationConflict
		message = "A Project already references this location"
		recoverable = true
	default:
		server.logger.Error("Project request failed", zap.Error(err))
	}

	return c.JSON(status, api.ErrorResponse{Error: api.Error{
		Code:        code,
		Message:     message,
		Recoverable: recoverable,
	}})
}

func cloneInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
