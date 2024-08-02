package controller

import (
	"github.com/jutimi/workspace-server/app/controller/api"
	"github.com/jutimi/workspace-server/app/controller/cms"
	"github.com/jutimi/workspace-server/app/middleware"
	"github.com/jutimi/workspace-server/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(
	router *gin.Engine,
	services service.ServiceCollections,
	middleware middleware.MiddlewareCollections,
) {
	api.NewApiWorkspaceController(router, services, middleware)

	cms.NewApiOrganizationController(router, services, middleware)
	cms.NewApiUserWorkspaceController(router, services, middleware)
	cms.NewApiWorkspaceController(router, services, middleware)
}
