package routers

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hexiaopi/rdm-toy/internal/app"
	"github.com/hexiaopi/rdm-toy/internal/middleware"
	v1 "github.com/hexiaopi/rdm-toy/internal/routers/api/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {

	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// profiling
	router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	// metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.Cors())
	apiRouter.Use(middleware.Logger())
	apiRouter.Use(middleware.Metrics())
	apiRouter.Use(middleware.Timeout(30 * time.Second))

	v1Router := apiRouter.Group("/v1")

	connController := v1.ConnController{}
	v1Router.GET("/conns", app.WrapList(connController.List))
	v1Router.GET("/conn/:conn", app.WrapData(connController.Get))
	v1Router.DELETE("/conn/:conn", app.Wrap(connController.Delete))

	v1Router.POST("/conn", app.Wrap(connController.Create))
	v1Router.POST("/conn/test", app.Wrap(connController.Test))
	v1Router.GET("/conn/:conn/clients", app.WrapData(connController.Clients))
	v1Router.DELETE("/conn/:conn/client", app.Wrap(connController.DeleteClient))

	dbController := v1.DBController{Client: connController}
	v1Router.GET("/conn/:conn/dbs", app.WrapList(dbController.List))

	keyController := v1.KeyController{Conn: dbController}
	v1Router.GET("/conn/:conn/db/:db/keys", app.WrapList(keyController.List))
	v1Router.DELETE("/conn/:conn/db/:db/keys", app.Wrap(keyController.DeleteKeys))
	v1Router.GET("/conn/:conn/db/:db/key/:key", app.WrapData(keyController.Get))
	v1Router.POST("/conn/:conn/db/:db/key", app.Wrap(keyController.Create))
	v1Router.DELETE("/conn/:conn/db/:db/key/:key", app.Wrap(keyController.Delete))
	v1Router.DELETE("/conn/:conn/db/:db/key/:key/field", app.Wrap(keyController.DeleteField))
	v1Router.PATCH("/conn/:conn/db/:db/key/:key/name", app.Wrap(keyController.PatchName))
	v1Router.PATCH("/conn/:conn/db/:db/key/:key/ttl", app.Wrap(keyController.PatchTTL))

	wsRouter := router.Group("/ws")
	wsRouter.GET("/conn/:conn/cmd", connController.Cmd)
	return router
}
