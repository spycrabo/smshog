package main

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"smshog/internal/http"
	ws "smshog/internal/websocket"
)

//go:embed assets/* templates/*
var embeds embed.FS

func main() {
	messages := make([]*http.Message, 0)

	hub := ws.NewHub()
	go hub.Run()

	gin.SetMode(gin.ReleaseMode)
	r := http.NewHttpServer(embeds, hub, messages)

	port := flag.String("p", "8080", "port to run at")
	flag.Parse()

	err := r.Run(":" + *port)
	if err != nil {
		panic(err)
	}
}
