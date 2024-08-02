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
)

type organizationHandler struct {
	services   service.ServiceCollections
	middleware middleware.MiddlewareCollections
}

func NewApiOrganizationController(
	router *gin.Engine,
	services service.ServiceCollections,
	middleware middleware.MiddlewareCollections,
) {
	handler := organizationHandler{services, middleware}

	group := router.Group("cms/v1/organizations", middleware.WorkspaceMW.Handler())
	{
		group.PUT("/update", handler.update)
		group.DELETE("/remove", handler.remove)
		group.POST("/create", handler.create)
	}
}

func (h *organizationHandler) update(c *gin.Context) {
	var data model.UpdateOrganizationRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.OrganizationSvc.UpdateOrganization(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *organizationHandler) create(c *gin.Context) {
	var data model.CreateOrganizationRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.OrganizationSvc.CreateOrganization(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *organizationHandler) remove(c *gin.Context) {
	var data model.RemoveOrganizationRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.OrganizationSvc.RemoveOrganization(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
