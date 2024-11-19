package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Player struct {
	ID   uuid.UUID `json:"id"`
	Num  int       `json:"num"`
	Name string    `json:"name"`
	ws   *websocket.Conn
}

type Game struct {
	players []*Player
}

func (game *Game) AppendPlayer(ws *websocket.Conn) *Player {
	player := Player{
		ID:   uuid.New(),
		Num:  0,
		Name: fmt.Sprintf("player %d:", len(game.players)),
		ws:   ws,
	}
	game.players = append(game.players, &player)
	return &player
}

func (game *Game) RemovePlayer(removePlayer *Player) {
	for i, player := range game.players {
		if player.ID == removePlayer.ID {
			game.players = append(game.players[:i], game.players[i+1:]...)
			break
		}
	}
}

type GameStateResource struct {
	Players       []*Player `json:"players"`
	CurrentPlayer *Player   `json:"current_player"`
}

func (game *Game) DispatchUpdate() {
	for _, player := range game.players {
		go func() {
			websocket.JSON.Send(player.ws, GameStateResource{
				Players:       game.players,
				CurrentPlayer: player,
			})
		}()
	}
}

func (game *Game) updatePlayer(ws *websocket.Conn) {
	player := game.AppendPlayer(ws)

	for {
		game.DispatchUpdate()

		buf := make([]byte, 1024)
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				game.RemovePlayer(player)
				game.DispatchUpdate()
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

		player.Num += int

		ws.Write([]byte(strconv.Itoa(player.Num)))
	}
}

func main() {
	game := Game{}

	http.HandleFunc("/", indexPage)
	http.Handle("/update-player", websocket.Handler(game.updatePlayer))

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
