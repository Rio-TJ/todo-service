package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	fmt.Println("Todo Service API v0.0")

	listenPort := flag.String("listenport", "4000", "Which port to listen")
	dbHost := flag.String("dbhost", "localhost", "PostgreSQL hostname")
	dbUser := flag.String("dbuser", "developer", "PostgreSQL user name")
	dbPassword := flag.String("dbpassword", "developer", "PostgreSQL user password")
	dbName := flag.String("dbname", "todo_service", "PostgreSQL database name")
	dbPort := flag.String("dbport", "5432", "PostgreSQL database port")
	frontendAppDomain := flag.String("frontendappdomain", "http://localhost:3000", "The domain of frontend application for allowing CORS requests")

	flag.Parse()

	db, err :=
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
