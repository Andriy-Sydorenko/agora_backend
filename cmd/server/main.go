package main

import (
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.Load("config.yml")
	router := gin.Default()
	router.SetTrustedProxies(cfg.Server.Cors.AllowedOrigins)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	if err := router.Run(); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
