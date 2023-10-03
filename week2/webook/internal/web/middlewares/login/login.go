// Copyright@daidai53 2023
package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MiddlewareBuilder struct {
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}

		sess := sessions.Default(context)
		if sess.Get("userId") == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
