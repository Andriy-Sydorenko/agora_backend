package main

import (
	"fmt"
	"log"

	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/Andriy-Sydorenko/agora_backend/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load("config.yml")

	mainRouter := router.SetupRouter(cfg)

	mainRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	if err := mainRouter.Run(fmt.Sprintf(":%d", cfg.Project.AppPort)); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
