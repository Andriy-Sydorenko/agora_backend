package main

import (
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Println(cfg)
	fmt.Println(cfg.JWT)
}
