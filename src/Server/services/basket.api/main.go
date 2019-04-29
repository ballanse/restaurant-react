package main

import (
	"fmt"
	"github.com/jurabek/basket.api/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jurabek/basket.api/middlewares"

	"github.com/jurabek/basket.api/controllers"
	"github.com/jurabek/basket.api/eureka"
	"github.com/jurabek/basket.api/repositories"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis" // swagger embed files

	// docs is generated by Swag CLI, you have to import it.
	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
)

// @title Basket API
// @version 1.0
// @description This is a rest api for basket which saves items to redis server
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securitydefinitions.oauth2.application Identity Server OAuth
// @tokenUrl http://localhost:5000/connect/token
func main() {
	os.Setenv("PORT", "5050")
	gin.SetMode(gin.DebugMode)

	handleSigterm()
	router := gin.Default()
	router.Use(middlewares.RequestMiddleware())

	redisPool, err := initRedis()

	if err != nil {
		fmt.Print(err)
	}

	connectionProvider := database.RedisConnectionProvider{
		Pool:redisPool,
	}

	basketRepository := repositories.NewRedisBasketRepository(&connectionProvider)
	controller := controllers.NewBasketController(basketRepository)

	auth := middlewares.CreateAuth()
	router.Use(auth.AuthMiddleware())

	api := router.Group("/api/v1/")
	{
		basket := api.Group("items")
		{
			basket.GET(":id", controller.Get)
			basket.POST("", controller.Create)
		}
	}

	// Home page should be redirected to swagger page
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	go eureka.Register()
	_ = router.Run()
}

func handleSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Kill, os.Interrupt)
	go func() {
		<-c
		eureka.UnRegister()
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}()
}

func initRedis() (*redis.Pool, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = ":6379"
	}
	pool := database.NewRedisPool(redisHost)
	database.CleanupHook(pool)

	err := database.HealthCheck(pool)

	return pool, err
}
