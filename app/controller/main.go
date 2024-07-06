package controller

import (
	"oauth-server/app/controller/api"
	"oauth-server/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(router *gin.Engine, services service.ServiceCollections) {
	api.NewApiOAuthController(router, services)
}
