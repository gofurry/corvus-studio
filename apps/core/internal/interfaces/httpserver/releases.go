package httpserver

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	releaseapp "github.com/gofurry/corvus-studio/apps/core/internal/application/release"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/httpserver/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ReleaseService interface {
	Create(context.Context, releaseapp.CreateCommand) (releasedomain.Goal, error)
	Get(context.Context, releasedomain.ID) (releasedomain.Goal, error)
	ListByProject(context.Context, projectdomain.ID) ([]releasedomain.Goal, error)
	Transition(context.Context, releasedomain.ID, releasedomain.Status) (releasedomain.Goal, error)
	LatestTemplate(string) (templatedomain.Definition, error)
}

func (server *Server) getReleaseTemplateHandler(c *echo.Context) error {
	definition, err := server.releases.LatestTemplate(c.Param("template_key"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := releaseTemplateResponse(definition)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (server *Server) createReleaseHandler(c *echo.Context) error {
	request, err := decodeWorkflowRequest[api.CreateReleaseRequest](c)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	projectID, err := projectdomain.IDFromUUID(request.ProjectId)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	goal, err := server.releases.Create(c.Request().Context(), releaseapp.CreateCommand{
		ProjectID: projectID, GoalType: releasedomain.GoalType(request.GoalType),
		TemplateKey: request.TemplateKey, TemplateVersion: request.TemplateVersion,
	})
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := releaseGoalResponse(goal)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	c.Response().Header().Set("Location", "/api/v1/releases/"+goal.ID.String())
	return c.JSON(http.StatusCreated, response)
}

func (server *Server) listProjectReleasesHandler(c *echo.Context) error {
	projectID, err := projectdomain.ParseID(c.Param("project_id"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	goals, err := server.releases.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	items := make([]api.ReleaseGoal, 0, len(goals))
	for _, goal := range goals {
		item, err := releaseGoalResponse(goal)
		if err != nil {
			return server.writeWorkflowError(c, err)
		}
		items = append(items, item)
	}
	return c.JSON(http.StatusOK, api.ReleaseGoalList{Items: items})
}

func (server *Server) getReleaseHandler(c *echo.Context) error {
	id, err := releasedomain.ParseID(c.Param("release_id"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	goal, err := server.releases.Get(c.Request().Context(), id)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := releaseGoalResponse(goal)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (server *Server) transitionReleaseHandler(c *echo.Context) error {
	id, err := releasedomain.ParseID(c.Param("release_id"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	request, err := decodeWorkflowRequest[api.TransitionReleaseRequest](c)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	goal, err := server.releases.Transition(c.Request().Context(), id, releasedomain.Status(request.Status))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := releaseGoalResponse(goal)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	return c.JSON(http.StatusOK, response)
}

func releaseGoalResponse(goal releasedomain.Goal) (api.ReleaseGoal, error) {
	id, err := uuid.Parse(goal.ID.String())
	if err != nil || id.Version() != 7 {
		return api.ReleaseGoal{}, fmt.Errorf("encode Release Goal ID %q: expected UUIDv7", goal.ID)
	}
	projectID, err := uuid.Parse(goal.ProjectID.String())
	if err != nil || projectID.Version() != 7 {
		return api.ReleaseGoal{}, fmt.Errorf("encode Release Project ID %q: expected UUIDv7", goal.ProjectID)
	}
	total, err := checkedInt32(goal.Summary.Total)
	if err != nil {
		return api.ReleaseGoal{}, err
	}
	done, err := checkedInt32(goal.Summary.Done)
	if err != nil {
		return api.ReleaseGoal{}, err
	}
	blocked, err := checkedInt32(goal.Summary.Blocked)
	if err != nil {
		return api.ReleaseGoal{}, err
	}
	requiredTotal, err := checkedInt32(goal.Summary.RequiredTotal)
	if err != nil {
		return api.ReleaseGoal{}, err
	}
	requiredDone, err := checkedInt32(goal.Summary.RequiredDone)
	if err != nil {
		return api.ReleaseGoal{}, err
	}
	return api.ReleaseGoal{
		Id: id, ProjectId: projectID, GoalType: api.ReleaseGoalType(goal.GoalType), Title: goal.Title,
		Status: api.ReleaseStatus(goal.Status), TemplateKey: goal.TemplateKey, TemplateVersion: goal.TemplateVersion,
		ChecklistSummary: api.ChecklistSummary{
			Total: total, Done: done, Blocked: blocked, RequiredTotal: requiredTotal, RequiredDone: requiredDone,
		},
		CreatedAt: goal.CreatedAt.UTC(), UpdatedAt: goal.UpdatedAt.UTC(),
	}, nil
}

func releaseTemplateResponse(definition templatedomain.Definition) (api.ReleaseTemplate, error) {
	reviewedAt, err := time.Parse(time.DateOnly, definition.ReviewedAt)
	if err != nil {
		return api.ReleaseTemplate{}, fmt.Errorf("encode template reviewed_at: %w", err)
	}
	items := make([]api.ReleaseTemplateItem, 0, len(definition.Items))
	for _, item := range definition.Items {
		sortOrder, err := checkedInt32(item.SortOrder)
		if err != nil {
			return api.ReleaseTemplate{}, err
		}
		items = append(items, api.ReleaseTemplateItem{
			TemplateItemKey: item.Key, Title: item.Title, Description: item.Description,
			Requirement: item.Requirement, Category: api.ChecklistCategory(item.Category),
			RequirementLevel: api.ChecklistRequirementLevel(item.RequirementLevel), Source: api.ChecklistSource(item.Source),
			SourceReference: item.SourceReference, SortOrder: sortOrder,
		})
	}
	return api.ReleaseTemplate{
		Key: definition.Key, Version: definition.Version, SchemaVersion: int32(definition.SchemaVersion),
		Name: definition.Name, Description: definition.Description, ReviewedAt: openapi_types.Date{Time: reviewedAt}, Items: items,
	}, nil
}

func checkedInt32(value int64) (int32, error) {
	if value < 0 || value > math.MaxInt32 {
		return 0, fmt.Errorf("value %d cannot be represented as int32", value)
	}
	return int32(value), nil
}
