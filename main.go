package main

import (
	"log"
	"platifyapi/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	server := gin.Default()
	routes.RegisterRoutes(server)
	if err := server.Run(":8080"); err != nil {
        log.Fatal("Unable to start server:", err)
    }
}