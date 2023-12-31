// Copyright@daidai53 2023
//go:build wireinject

package main

import (
	"github.com/daidai53/webook/internal/repository"
	"github.com/daidai53/webook/internal/repository/cache"
	"github.com/daidai53/webook/internal/repository/dao"
	"github.com/daidai53/webook/internal/service"
	"github.com/daidai53/webook/internal/web"
	"github.com/daidai53/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB,
		ioc.InitRedisClient,
		dao.NewUserDAO,
		ioc.NewLocalCacheDefault,

		// cache部分
		//cache.NewRedisCodeCache,
		cache.NewLocalCodeCache,
		cache.NewUserCache,

		// repository部分
		repository.NewCachedCodeRepository,
		repository.NewCachedUserRepository,

		// service部分
		ioc.InitSmsService,
		service.NewUserService,
		service.NewCodeService,

		// handler部分
		web.NewUserHandler,

		ioc.InitWebServer,
		ioc.InitGinMiddlewares,
	)
	return gin.Default()
}
