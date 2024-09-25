// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"geektime/webook/internal/repository"
	"geektime/webook/internal/repository/cache"
	"geektime/webook/internal/repository/dao"
	"geektime/webook/internal/service"
	"geektime/webook/internal/web"
	"geektime/webook/internal/web/jwt"
	"geektime/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitArticleHandler(dao2 dao.ArticleDAO) *web.ArticleHandler {
	cmdable := InitRedis()
	articleCache := cache.NewArticleRedisCache(cmdable)
	loggerV1 := InitLogger()
	db := InitDB()
	userDAO := dao.NewUserDao(db)
	articleRepository := repository.NewArticleRepository(dao2, articleCache, loggerV1, userDAO)
	articleService := service.NewArticleService(articleRepository, loggerV1)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, loggerV1, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1, interactiveService)
	return articleHandler
}

func InitInteractiveService() service.InteractiveService {
	db := InitDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	loggerV1 := InitLogger()
	cmdable := InitRedis()
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, loggerV1, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	return interactiveService
}

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	jwtHandler := jwt.NewRedisJWTHandler(cmdable)
	loggerV1 := ioc.InitLoggerV1()
	v := ioc.InitMiddlewares(jwtHandler, loggerV1)
	db := ioc.InitDB(loggerV1)
	userDAO := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, jwtHandler)
	wechatService := ioc.InitWechatService()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, jwtHandler)
	articleDAO := dao.NewGROMArticleDAO(db)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewArticleRepository(articleDAO, articleCache, loggerV1, userDAO)
	articleService := service.NewArticleService(articleRepository, loggerV1)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, loggerV1, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1, interactiveService)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	return engine
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitRedis, InitDB,
	InitLogger)

var userSvcProvider = wire.NewSet(dao.NewUserDao, cache.NewUserCache, repository.NewUserRepository, service.NewUserService)

var articlSvcProvider = wire.NewSet(repository.NewArticleRepository, cache.NewArticleRedisCache, dao.NewGROMArticleDAO, service.NewArticleService)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)
