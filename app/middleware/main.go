package middleware

import (
	"net/http"
	"workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/gin-gonic/gin"
)

type userMiddleware struct{}

func NewUserMiddleware() Middleware {
	return &userMiddleware{}
}

func (m *userMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		userToken := c.Request.Header.Get(utils.USER_ACCESS_TOKEN_KEY)
		if userToken == "" {
			err := errors.New(errors.ErrCodeUnauthorized)
			c.JSON(http.StatusUnauthorized, utils.FormatErrorResponse(err))
			return
		}

		c.Next()
	}
}
