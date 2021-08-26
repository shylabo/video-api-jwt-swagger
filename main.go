package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shylabo/golang-gin-poc/api"
	"github.com/shylabo/golang-gin-poc/controller"
	"github.com/shylabo/golang-gin-poc/docs"
	"github.com/shylabo/golang-gin-poc/middlewares"
	"github.com/shylabo/golang-gin-poc/repository"
	"github.com/shylabo/golang-gin-poc/service"
	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
)

var (
	videoRepository repository.VideoRepository = repository.NewVideoRepository()
	videoService    service.VideoService       = service.New(videoRepository)
	loginService    service.LoginService       = service.NewLoginService()
	jwtService      service.JWTService         = service.NewJWTService()

	videoController controller.VideoController = controller.New(videoService)
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}
func main() {

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Golang-Gin-Poc"
	docs.SwaggerInfo.Description = "Youtube Video API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:5000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	defer videoRepository.CloseDB()

	setupLogOutput()
	server := gin.New()

	server.Use(gin.Recovery(), middlewares.Logger())

	server.Static("/css", "./templates/css")

	server.LoadHTMLGlob("templates/*.html")

	videoAPI := api.NewVideoAPI(loginController, videoController)

	apiRoutes := server.Group(docs.SwaggerInfo.BasePath)
	{
		login := apiRoutes.Group("/auth")
		{
			login.POST("/token", videoAPI.Authenticate)
		}

		// TODO: Fix JWT error in Swagger
		// videos := apiRoutes.Group("/videos", middlewares.AuthorizeJWT())
		videos := apiRoutes.Group("/videos")
		{
			videos.GET("", videoAPI.GetVideos)
			videos.POST("", videoAPI.CreateVideo)
			videos.PUT(":id", videoAPI.UpdateVideo)
			videos.DELETE(":id", videoAPI.DeleteVideo)
		}

		// TODO: Fix JWT error in Swagger
		// views := apiRoutes.Group("/view", middlewares.AuthorizeJWT())
		views := apiRoutes.Group("/view")
		{
			views.GET("/videos", videoController.ShowAll)
		}
	}

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// We can setup this env variable from the EB console
	port := os.Getenv("PORT")
	// Elastic Beanstalk forwards requests to port 5000
	if port == "" {
		port = "5000"
	}

	server.Run(":" + port)
}
