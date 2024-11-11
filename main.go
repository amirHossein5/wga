package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns []*websocket.Conn
}

func NewServer() *Server {
	return &Server{
		conns: make([]*websocket.Conn, 0),
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", indexPage)

	port := 8080
	log.Println("starting server in port", port)
	err := http.ListenAndServe(fmt.Sprint(":", port), nil)
	if err != nil {
		log.Fatal("failed to start server", err)
	}
}
