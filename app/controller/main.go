package controller

import (
	"workspace-server/app/controller/api"
	"workspace-server/app/controller/cms"
	"workspace-server/app/middleware"
	"workspace-server/app/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func RegisterControllers(
	router *gin.Engine,
	tracer trace.Tracer,
	middleware middleware.MiddlewareCollections,
	services service.ServiceCollections,
) {
	api.NewApiWorkspaceController(router, tracer, middleware, services)

	cms.NewApiOrganizationController(router, tracer, middleware, services)
	cms.NewApiUserWorkspaceController(router, tracer, middleware, services)
	cms.NewApiWorkspaceController(router, tracer, middleware, services)
}
