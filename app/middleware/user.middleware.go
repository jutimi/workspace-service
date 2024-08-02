package middleware

import (
	"net/http"

	"github.com/jutimi/workspace-server/package/errors"
	"github.com/jutimi/workspace-server/utils"

	"strings"

	"github.com/gin-gonic/gin"
)

type userMiddleware struct {
}

func NewUserMiddleware() Middleware {
	return &userMiddleware{}
}
func (m *userMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resErr := errors.New(errors.ErrCodeUnauthorized)

		token := c.GetHeader(utils.USER_AUTHORIZATION)
		if token == "" {
			c.JSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}
		tokenArr := strings.Split(token, " ")
		if len(tokenArr) != 2 {
			c.JSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}
		if tokenArr[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}

		payload, err := utils.ParseUserToken(tokenArr[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}

		c.Set(string(utils.USER_CONTEXT_KEY), payload)

		c.Next()
	}
}
