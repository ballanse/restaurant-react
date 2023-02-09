package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jurabek/basket.api/cmd/config"
	"github.com/jurabek/basket.api/internal/database"
	"github.com/jurabek/basket.api/internal/docs"
	"github.com/jurabek/basket.api/internal/handlers"
	"github.com/jurabek/basket.api/internal/middlewares"
	"github.com/jurabek/basket.api/pkg/producer"

	"github.com/jurabek/basket.api/internal/repositories"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis" // swagger embed files

	// docs is generated by Swag CLI, you have to import it.
	ginSwagger "github.com/Jurabek/gin-swagger" // gin-swagger middleware
	swaggerFiles "github.com/swaggo/files"      // swagger embed files
)

var (
	GitCommit string
	Version   string
)

//	@title			Basket API
//	@version		1.0
//	@description	This is a rest api for basket which saves items to redis server
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	gin.SetMode(gin.DebugMode)

	basePath, _ := os.LookupEnv("BASE_PATH")
	docs.SwaggerInfo.BasePath = basePath

	handleSigterm()
	router := gin.Default()
	router.Use(middlewares.RequestMiddleware())

	cfg := config.Init()

	redisPool, err := initRedis(cfg.RedisHost)
	if err != nil {
		fmt.Print(err)
	}

	p, err := sarama.NewSyncProducer([]string{cfg.KafkaBroker}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("new producer failed!")
	}
	defer p.Close()

	kafkaProducer := producer.NewKafkaEventProducer(p, cfg.CheckoutTopic)
	connectionProvider := database.RedisConnectionProvider{
		Pool: redisPool,
	}

	basketRepository := repositories.NewRedisBasketRepository(&connectionProvider)
	basketHandler := handlers.NewBasketHandler(basketRepository)
	checkoutHandler := handlers.NewCheckOutHandler(basketRepository, kafkaProducer)

	api := router.Group(basePath + "/api/v1/")
	{
		basket := api.Group("items")
		{
			basket.GET(":id", basketHandler.Get)
			basket.POST("", basketHandler.Create)
			basket.DELETE(":id", basketHandler.Delete)
		}

		checkout := api.Group("checkout")
		{
			checkout.POST("", checkoutHandler.Checkout)
		}
	}

	// Home page should be redirected to swagger page
	router.GET(basePath+"/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, basePath+"/swagger/index.html")
	})

	router.GET(basePath+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		func(c *ginSwagger.Config) {
			c.URL = basePath + "/swagger/doc.json"
		}))

	_ = router.Run()
}

func handleSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Kill, os.Interrupt)
	go func() {
		<-c
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}()
}

func initRedis(redisHost string) (*redis.Pool, error) {
	if redisHost == "" {
		redisHost = ":6379"
	}
	pool := database.NewRedisPool(redisHost)
	database.CleanupHook(pool)

	err := database.HealthCheck(pool)

	return pool, err
}
