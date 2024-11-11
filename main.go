package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

type Player struct {
	num int
}

var player Player = Player{num: 0}

func updatePlayer(ws *websocket.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("failed to read websocket data", err)
			continue
		}

		msg := string(buf[:n])

		int, err := strconv.Atoi(msg)
		if err != nil {
			log.Println(msg, "is not a number")
			continue
		}

		player.num += int

		ws.Write([]byte(strconv.Itoa(player.num)))
	}
}

func main() {
	http.HandleFunc("/", indexPage)
	http.Handle("/update-player", websocket.Handler(updatePlayer))

	port := 8080
	log.Println("starting server in port", port)
	err := http.ListenAndServe(fmt.Sprint(":", port), nil)
	if err != nil {
		log.Fatal("failed to start server", err)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
