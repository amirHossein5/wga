package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/amirhossein5/wgo/pkg/rand"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Player struct {
	ID        uuid.UUID `json:"id"`
	PositionY int       `json:"position_y"`
	PositionX int       `json:"position_x"`
	Color     string    `json:"color"`
	ws        *websocket.Conn
}

type Game struct {
	players []*Player
}

const (
	MOVE_SPEED = 30
	TO_RIGHT   = "to-right"
	TO_LEFT    = "to-left"
	TO_TOP     = "to-top"
	TO_DOWN    = "to-down"
)

func (game *Game) IsValidMove(move string) bool {
	return slices.Contains([]string{TO_RIGHT, TO_LEFT, TO_DOWN, TO_TOP}, move)
}

func (game *Game) HandleMove(player *Player, move string) error {
	if !game.IsValidMove(move) {
		return fmt.Errorf("not a valid move %v\n", move)
	}

	if move == TO_RIGHT {
		player.PositionX += MOVE_SPEED
	} else if move == TO_LEFT {
		player.PositionX -= MOVE_SPEED
	} else if move == TO_TOP {
		player.PositionY -= MOVE_SPEED
	} else if move == TO_DOWN {
		player.PositionY += MOVE_SPEED
	}

	return nil
}

func (game *Game) AppendPlayer(ws *websocket.Conn) *Player {
	player := Player{
		ID:        uuid.New(),
		PositionY: 200,
		PositionX: 500,
		Color:     rand.HexColor(),
		ws:        ws,
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

	loop(30, func() bool {
		game.DispatchUpdate()

		buf := make([]byte, 1024)
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				game.RemovePlayer(player)
				game.DispatchUpdate()
				return false
			}
			log.Println("failed to read websocket data", err)
			return true
		}

		err = game.HandleMove(player, string(buf[:n]))
		if err != nil {
			log.Printf("counldn't handle move %v\n", err)
			return true
		}

		return true
	})
}

func main() {
	game := Game{}

	http.HandleFunc("/", indexPage)
	http.Handle("/update-player", websocket.Handler(game.updatePlayer))

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

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

func loop(fps int, callback func() bool) {
	lastFrameTime := time.Now()
	delayPerFrame := time.Millisecond * time.Duration(1000/fps)

	for {
		currentTime := time.Now()
		elapsedTime := currentTime.Sub(lastFrameTime)
		if elapsedTime >= delayPerFrame {
			res := callback()
			if res == false {
				break
			}
			lastFrameTime = currentTime
		}

		time.Sleep(delayPerFrame - elapsedTime)
	}
}
