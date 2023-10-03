// Copyright@daidai53 2023
package main

import (
	"github.com/daidai53/go-homework/week2/webook/internal/repository"
	"github.com/daidai53/go-homework/week2/webook/internal/repository/dao"
	"github.com/daidai53/go-homework/week2/webook/internal/service"
	"github.com/daidai53/go-homework/week2/webook/internal/web"
	"github.com/daidai53/go-homework/week2/webook/internal/web/middlewares/login"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	initUserHdl(db, server)
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	handler := web.NewUserHandler(us)
	handler.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:12345)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initServer() *gin.Engine {
	server := gin.Default()
	login := login.MiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret"))
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders: []string{
			"Content-Type",
			"authorization",
		},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}), sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
