package api

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

type workspaceHandler struct {
	tracer     trace.Tracer
	middleware middleware.MiddlewareCollections
	services   service.ServiceCollections
}

func NewApiWorkspaceController(
	router *gin.Engine,
	tracer trace.Tracer,
	middleware middleware.MiddlewareCollections,
	services service.ServiceCollections,
) {
	handler := workspaceHandler{tracer, middleware, services}

	group := router.Group("api/v1/workspaces")
	{
		group.POST("/create", middleware.UserMW.Handler(), handler.create)
	}
}

func (h *workspaceHandler) create(c *gin.Context) {
	var data model.CreateWorkspaceRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx, main := h.tracer.Start(ctx, "create-workspace")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.WorkspaceSvc.CreateWorkspace(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
