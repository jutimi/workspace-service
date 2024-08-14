package cms

import (
	"context"
	"net/http"
	"time"

	"workspace-server/app/middleware"
	"workspace-server/app/model"
	"workspace-server/app/service"
	_errors "workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/trace"
)

type userWorkspaceHandler struct {
	tracer     trace.Tracer
	middleware middleware.MiddlewareCollections
	services   service.ServiceCollections
}

func NewApiUserWorkspaceController(
	router *gin.Engine,
	tracer trace.Tracer,
	middleware middleware.MiddlewareCollections,
	services service.ServiceCollections,
) {
	handler := userWorkspaceHandler{tracer, middleware, services}

	group := router.Group("cms/v1/user-workspaces", middleware.WorkspaceMW.Handler())
	{
		group.POST("/create", handler.create)
		group.PUT("/update", handler.update)
	}
}

func (h *userWorkspaceHandler) create(c *gin.Context) {
	var data model.CreateUserWorkspaceRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx, main := h.tracer.Start(ctx, "create-user-workspace")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.UserWorkspaceSvc.CreateUserWorkspace(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *userWorkspaceHandler) update(c *gin.Context) {
	var data model.UpdateUserWorkspaceRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx, main := h.tracer.Start(ctx, "update-user-workspace")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.UserWorkspaceSvc.UpdateUserWorkspace(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
