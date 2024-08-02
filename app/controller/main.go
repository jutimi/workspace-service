package controller

import (
	"workspace-server/app/controller/api"
	"workspace-server/app/controller/cms"
	"workspace-server/app/middleware"
	"workspace-server/app/service"

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
