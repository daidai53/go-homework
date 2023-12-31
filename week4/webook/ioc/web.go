// Copyright@daidai53 2023
package ioc

import (
	"github.com/daidai53/webook/internal/web"
	"github.com/daidai53/webook/internal/web/middlewares/login"
	"github.com/daidai53/webook/pkg/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdlw []gin.HandlerFunc, handlers *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdlw...)
	handlers.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient goredis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders: []string{
				"Content-Type",
				"authorization",
			},
			ExposeHeaders: []string{
				"x-jwt-token",
			},
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 12 * time.Hour,
		}),
		(&login.MiddlewareJWTBuilder{}).CheckLogin(),
		ratelimit.NewRedisSliceWindowLimiter(redisClient, 30*time.Second, 1000).BuildLua(),
	}
}

//
//func useSession(server *gin.Engine) {
//	login := login.MiddlewareBuilder{}
//	//store := memstore.NewStore([]byte("jfdaflalkfhlaf"), []byte("ffadfadfsadfasad"))
//	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
//		[]byte("jfdaflalkfhlaf"), []byte("ffadfadfsadfasad"))
//	if err != nil {
//		panic(err)
//	}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
