package main

// @title           Collector Ouphe API
// @version         1.0
// @description     Сервис сбора и анализа данных Collector Ouphe
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
//
// @host      localhost:8080
// @BasePath  /

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ShenokZlob/collector-service/docs"

	"github.com/ShenokZlob/collector-service/internal/controllers"
	"github.com/ShenokZlob/collector-service/internal/controllers/middleware"
	repositories "github.com/ShenokZlob/collector-service/internal/rep/mongo"
	"github.com/ShenokZlob/collector-service/pkg/logger"
	"github.com/ShenokZlob/collector-service/usecase/auth"
	"github.com/ShenokZlob/collector-service/usecase/collection"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Starting collector service")

	if os.Getenv("JWT_SECRET") == "" {
		panic("Don't have JWT_SECRET")
	}

	// Init db
	log.Info("Init database")
	connString := os.Getenv("MONGO_CONN_STRING")
	db, err := mongo.Connect(options.Client().ApplyURI(connString))
	if err != nil {
		panic(err)
	}
	if err := db.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	// Init server
	log.Info("Init app server")
	host := os.Getenv("SERVER_ADDRESS")

	rep := repositories.NewRepository(db)

	servAuth := auth.NewAuthUsecase(log, rep)
	servCollections := collection.NewCollectionsService(log, rep)
	servCards := collection.NewCardsService(log, rep)

	ctrlAuth := controllers.NewAuthController(log, servAuth)
	ctrlCollections := controllers.NewCollectionsController(log, servCollections)
	ctrlCards := controllers.NewCardsController(log, servCards)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Public routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/")
	{
		public.POST("/register", ctrlAuth.Register)
		public.POST("/login", ctrlAuth.Login)
		public.POST("/refresh", ctrlAuth.RefreshToken)
		public.POST("/logout", ctrlAuth.Logout)
	}

	publicTelegram := router.Group("/telegram")
	{
		publicTelegram.POST("/register", ctrlAuth.RegisterTelegram)
		publicTelegram.GET("/link", ctrlAuth.LinkTelegram)
		publicTelegram.POST("/unlink", ctrlAuth.UnlinkTelegram)
	}

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(log)
	authorized := router.Group("/", authMiddleware)
	{
		authorized.GET("/collections", ctrlCollections.GetAll)
		authorized.GET("/collections/:id", ctrlCollections.Get)
		authorized.POST("/collections", ctrlCollections.Create)
		authorized.PATCH("/collections/:id", ctrlCollections.Rename)
		authorized.DELETE("/collections/:id", ctrlCollections.Delete)
		// authorized.GET("/collections/name/:name", ctrlCollections.GetCollectionByName)

		authorized.GET("/collections/:id/cards", ctrlCards.ListCardsInCollection)
		authorized.POST("/collections/:id/cards", ctrlCards.AddCardToCollection)
		authorized.PATCH("/collections/:id/cards/:id", ctrlCards.SetCardCountInCollection)
		authorized.DELETE("/collections/:id/cards/:id", ctrlCards.DeleteCardFromCollection)
	}

	server := &http.Server{
		Addr:    host,
		Handler: router.Handler(),
	}

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log.Info("Starting server", zap.String("host", host))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("ListenAndServe", zap.Error(err))
		}
	}()
	log.Info("Server started", zap.String("host", host))

	// Stop server
	<-ctx.Done()

	log.Info("Stopping app")
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown Failed", zap.Error(err))
	}

	os.Exit(0)
}
