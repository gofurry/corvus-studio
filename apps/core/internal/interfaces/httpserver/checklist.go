package httpserver

import (
	"context"
	"fmt"
	"net/http"

	checklistapp "github.com/gofurry/corvus-studio/apps/core/internal/application/checklist"
	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/httpserver/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type ChecklistService interface {
	Create(context.Context, checklistapp.CreateCommand) (checklistdomain.Item, error)
	Get(context.Context, checklistdomain.ID) (checklistdomain.Item, error)
	List(context.Context, releasedomain.ID, checklistdomain.Filter) ([]checklistdomain.Item, error)
	Transition(context.Context, checklistdomain.ID, checklistdomain.Status) (checklistdomain.Item, error)
}

func (server *Server) createChecklistItemHandler(c *echo.Context) error {
	request, err := decodeWorkflowRequest[api.CreateChecklistItemRequest](c)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	releaseID, err := releasedomain.IDFromUUID(request.ReleaseGoalId)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	item, err := server.checklist.Create(c.Request().Context(), checklistapp.CreateCommand{
		ReleaseGoalID: releaseID, Title: request.Title, Description: request.Description,
		Requirement: request.Requirement, Category: checklistdomain.Category(request.Category),
		RequirementLevel: checklistdomain.RequirementLevel(request.RequirementLevel),
	})
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := checklistItemResponse(item)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	c.Response().Header().Set("Location", "/api/v1/checklist/"+item.ID.String())
	return c.JSON(http.StatusCreated, response)
}

func (server *Server) getChecklistItemHandler(c *echo.Context) error {
	id, err := checklistdomain.ParseID(c.Param("item_id"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	item, err := server.checklist.Get(c.Request().Context(), id)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := checklistItemResponse(item)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (server *Server) listChecklistItemsHandler(c *echo.Context) error {
	releaseID, err := releasedomain.ParseID(c.Param("release_id"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	filter := checklistdomain.Filter{}
	if value := c.QueryParam("status"); value != "" {
		status := checklistdomain.Status(value)
		filter.Status = &status
	}
	if value := c.QueryParam("source"); value != "" {
		source := checklistdomain.Source(value)
		filter.Source = &source
	}
	if value := c.QueryParam("category"); value != "" {
		category := checklistdomain.Category(value)
		filter.Category = &category
	}
	items, err := server.checklist.List(c.Request().Context(), releaseID, filter)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	responseItems := make([]api.ChecklistItem, 0, len(items))
	for _, item := range items {
		response, err := checklistItemResponse(item)
		if err != nil {
			return server.writeWorkflowError(c, err)
		}
		responseItems = append(responseItems, response)
	}
	return c.JSON(http.StatusOK, api.ChecklistItemList{Items: responseItems})
}

func (server *Server) transitionChecklistItemHandler(c *echo.Context) error {
	id, err := checklistdomain.ParseID(c.Param("item_id"))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	request, err := decodeWorkflowRequest[api.TransitionChecklistItemRequest](c)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	item, err := server.checklist.Transition(c.Request().Context(), id, checklistdomain.Status(request.Status))
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	response, err := checklistItemResponse(item)
	if err != nil {
		return server.writeWorkflowError(c, err)
	}
	return c.JSON(http.StatusOK, response)
}

func checklistItemResponse(item checklistdomain.Item) (api.ChecklistItem, error) {
	id, err := uuid.Parse(item.ID.String())
	if err != nil || id.Version() != 7 {
		return api.ChecklistItem{}, fmt.Errorf("encode Checklist item ID %q: expected UUIDv7", item.ID)
	}
	releaseID, err := uuid.Parse(item.ReleaseGoalID.String())
	if err != nil || releaseID.Version() != 7 {
		return api.ChecklistItem{}, fmt.Errorf("encode Checklist Release ID %q: expected UUIDv7", item.ReleaseGoalID)
	}
	sortOrder, err := checkedInt32(item.SortOrder)
	if err != nil {
		return api.ChecklistItem{}, err
	}
	return api.ChecklistItem{
		Id: id, ReleaseGoalId: releaseID, Title: item.Title, Description: item.Description,
		Requirement: item.Requirement, Category: api.ChecklistCategory(item.Category),
		RequirementLevel: api.ChecklistRequirementLevel(item.RequirementLevel), Source: api.ChecklistSource(item.Source),
		SourceReference: item.SourceReference, TemplateItemKey: cloneString(item.TemplateItemKey),
		TemplateKey: cloneString(item.TemplateKey), TemplateVersion: cloneString(item.TemplateVersion),
		Status: api.ChecklistStatus(item.Status), SortOrder: sortOrder,
		CreatedAt: item.CreatedAt.UTC(), UpdatedAt: item.UpdatedAt.UTC(),
	}, nil
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
