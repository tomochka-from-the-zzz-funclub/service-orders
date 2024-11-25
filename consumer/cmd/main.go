package main

import (
	"consumer/internal/config"
	"consumer/internal/service"
	"consumer/internal/transport"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("%v", cfg)
	s := service.NewSrv(cfg)
	go s.Read(cfg)
	transport.HandleCreate(cfg, s)
}
