package main

import (
	"flag"
	"fmt"
	"github.com/SoftclubIT/todo-service/pkg/database"
	"github.com/SoftclubIT/todo-service/pkg/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println("Todo Service API v0.0")

	listenPort := flag.String("listenport", "4000", "Which port to listen")
	dbHost := flag.String("dbhost", "localhost", "PostgreSQL hostname")
	dbUser := flag.String("dbuser", "developer", "PostgreSQL user name")
	dbPassword := flag.String("dbpassword", "developer", "PostgreSQL user password")
	dbName := flag.String("dbname", "todo_service", "PostgreSQL database name")
	dbPort := flag.String("dbport", "5432", "PostgreSQL database port")
	//frontendAppDomain := flag.String("frontendappdomain", "http://localhost:3000", "The domain of frontend application for allowing CORS requests")

	flag.Parse()

	db, err := database.Init(*dbHost, *dbUser, *dbPassword, *dbName, *dbPort)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
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

	router.Run(":" + *listenPort)
}
