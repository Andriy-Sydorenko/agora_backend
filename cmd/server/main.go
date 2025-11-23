package main

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.Load("config.yml")

	mainRouter := router.SetupRouter(cfg)

	mainRouter.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	if err := mainRouter.Run(); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
