package middleware

import (
	"net/http"

	"workspace-server/package/errors"
	"workspace-server/utils"

	"strings"

	"github.com/gin-gonic/gin"
)

type workspaceMiddleware struct {
}

func NewWorkspaceMiddleware() Middleware {
	return &workspaceMiddleware{}
}
func (m *workspaceMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resErr := errors.New(errors.ErrCodeUnauthorized)

		token := c.GetHeader(utils.USER_AUTHORIZATION)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}
		tokenArr := strings.Split(token, " ")
		if len(tokenArr) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}
		if tokenArr[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}

		payload, err := utils.ParseWorkspaceToken(tokenArr[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}

		c.Set(string(utils.WORKSPACE_CONTEXT_KEY), payload)

		c.Next()
	}
}
