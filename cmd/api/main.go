package main

import (
	"flag"
	"fmt"
	"github.com/SoftclubIT/todo-service/docs"
	_ "github.com/SoftclubIT/todo-service/docs"
	"github.com/SoftclubIT/todo-service/pkg/database"
	"github.com/SoftclubIT/todo-service/pkg/handlers"
	"github.com/SoftclubIT/todo-service/pkg/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

//	@title			Todo Service API
//	@version		0.0.0
//	@description	This is a sample Todo Application API

//	@contact.name	SoftClub
//	@contact.url	https://www.softclub.tj/Contacts
//	@contact.email	info@softclub.tj

// @host						localhost:4000
// @BasePath					/
// @query.collection.format	multi
func main() {
	fmt.Println("Todo Service API v0.0")

	host := flag.String("host", "http://localhost", "The server address")
	port := flag.String("listenport", "4000", "Which port to listen")
	dbHost := flag.String("dbhost", "localhost", "PostgreSQL hostname")
	dbUser := flag.String("dbuser", "developer", "PostgreSQL user name")
	dbPassword := flag.String("dbpassword", "developer", "PostgreSQL user password")
	dbName := flag.String("dbname", "todo_service", "PostgreSQL database name")
	dbPort := flag.String("dbport", "5432", "PostgreSQL database port")
	//frontendAppDomain := flag.String("frontendappdomain", "http://localhost:3000", "The domain of frontend application for allowing CORS requests")

	flag.Parse()

	docs.SwaggerInfo.Host = *host
	docs.SwaggerInfo.Schemes = []string{"http"}

	db, err := database.Init(*dbHost, *dbUser, *dbPassword, *dbName, *dbPort)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware(""))
	hndlrs := handlers.New(db)

	router.GET("/health-check", hndlrs.HealthCheck)

	tasks := router.Group("tasks")
	{
		tasks.GET("", hndlrs.GetTasks)
		tasks.POST("", hndlrs.CreateTask)
		tasks.DELETE("/:taskID", hndlrs.DeleteTask)
		tasks.POST("/:taskID/completed", hndlrs.CompleteTask)
		tasks.DELETE("/:taskID/completed", hndlrs.UndoTask)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + *port)
}
