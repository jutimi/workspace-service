package controller

import (
	"workspace-server/app/controller/api"
	"workspace-server/app/controller/cms"
	"workspace-server/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(router *gin.Engine, services service.ServiceCollections) {
	api.NewApiWorkspaceController(router, services)

	cms.NewApiOrganizationController(router, services)
	cms.NewApiUserWorkspaceController(router, services)
	cms.NewApiWorkspaceController(router, services)
}
