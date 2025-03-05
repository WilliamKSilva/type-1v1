package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/WilliamKSilva/type-1v1/internal"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: Mensagem de connection tem que ser a primeira mensagem obrigatoriamente

type State struct {
	Mu    sync.Mutex
	Rooms *[]internal.Room
}

var upgrader = websocket.Upgrader{}

func unmarshalMessage[T any](c *websocket.Conn, mt int, buf []byte) *T {
	var message T
	err := json.Unmarshal(buf, &message)
	if err != nil {
		log.Println("Error trying to parse message")
		err = c.WriteMessage(mt, []byte("Invalid message format"))
		if err != nil {
			log.Println("Error trying to send warning message")
			return nil
		}

		return nil
	}

	return &message
}

func broadcastMessages(c *websocket.Conn, conn internal.Connection, state *State) {
	defer c.Close()

	for {
		state.Mu.Lock()
		for i, r := range *state.Rooms {
			if conn.RoomId != r.Id {
				continue
			}

			mes := &(*state.Rooms)[i].Messages
			for j, m := range *mes {
				if m.DeliveredCount >= len(r.Connections)-1 {
					(*mes) = slices.Delete((*mes), j, j+1)
					continue
				}

				if m.DeliveredTo[conn.Id] || m.ConnectionId == conn.Id {
					continue
				}

				data, err := json.Marshal(m)
				if err != nil {
					log.Println("[broadcastMessages] error trying to marshal message:", err)
				}
				c.WriteMessage(websocket.BinaryMessage, data)
				(*mes)[j].DeliveredTo[conn.Id] = true
				(*mes)[j].DeliveredCount += 1
			}

			continue
		}
		state.Mu.Unlock()
	}
}

func readMessages(c *websocket.Conn, conn internal.Connection, state *State) {
	defer c.Close()

	for {
		mt, buf, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		data := unmarshalMessage[internal.Message](c, mt, buf)
		if data == nil {
			continue
		}

		state.Mu.Lock()
		for i, r := range *state.Rooms {
			if data.RoomId == r.Id {
				(*state.Rooms)[i].Messages = append((*state.Rooms)[i].Messages, internal.Message{
					Id:           uuid.New().String(),
					Content:      data.Content,
					RoomId:       data.RoomId,
					ConnectionId: conn.Id,
					DeliveredTo:  make(map[string]bool),
				})
			}

			log.Printf("New message received from: %s. Content: %s", conn.Name, data.Content)
		}
		state.Mu.Unlock()
	}
}

func socket_conn(w http.ResponseWriter, r *http.Request, state *State) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	// defer c.Close()

	mt, buf, err := c.ReadMessage()
	if err != nil {
		log.Println("Error trying to read connection message")
		return
	}

	data := unmarshalMessage[internal.Connection](c, mt, buf)
	if data == nil {
		return
	}

	connId := uuid.New().String()
	conn := internal.Connection{
		Id:     connId,
		Name:   data.Name,
		RoomId: data.RoomId,
	}

	for i, r := range *state.Rooms {
		if r.Id == conn.RoomId {
			(*state.Rooms)[i].Connections = append((*state.Rooms)[i].Connections, conn)
		}
	}

	res, err := json.Marshal(data)
	if err != nil {
		log.Println("error: ", err)
		return
	}
	c.WriteMessage(mt, res)
	log.Printf("Connected: %s with Id: %s", conn.Name, conn.Id)

	go readMessages(c, conn, state)
	go broadcastMessages(c, conn, state)
}

const port = "8080"

func main() {
	flag.Parse()

	const mockedRoomId string = "21ca15d0-e346-4630-a240-773a828c31b3"

	state := State{
		Rooms: &[]internal.Room{
			{
				Id: mockedRoomId,
			},
		},
	}

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) { socket_conn(w, r, &state) })
	log.Printf("Server running at port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), nil))
}
