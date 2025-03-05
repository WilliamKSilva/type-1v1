package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/WilliamKSilva/type-1v1/internal"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: Mensagem de connection tem que ser a primeira mensagem obrigatoriamente

var upgrader = websocket.Upgrader{}

func unmarshal_message[T any](c *websocket.Conn, mt int, buf []byte) *T {
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

func broadcast_messages(c *websocket.Conn, conn internal.Connection, rooms *[]internal.Room) {
	for {
		for i, r := range *rooms {
			if conn.RoomId != r.Id {
				continue
			}

			// TODO: remover mensagem da slice quando já tiver sido entregue
			// para todas as conexões do Room
			for j, m := range (*rooms)[i].Messages {
				// Mensagem já foi entregue para essa conexão
				// log.Println(m.DeliveredTo[connId])
				if m.DeliveredTo[conn.Id] || m.ConnectionId == conn.Id {
					continue
				}

				c.WriteMessage(1, []byte(m.Content))
				(*rooms)[i].Messages[j].DeliveredTo[conn.Id] = true
			}

			continue
		}
	}
}

func read_messages(c *websocket.Conn, conn internal.Connection, rooms *[]internal.Room) {
	// TODO: tratar o fechamento da conexão
	for {
		mt, buf, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		data := unmarshal_message[internal.Message](c, mt, buf)
		if data == nil {
			continue
		}

		for i, r := range *rooms {
			if data.RoomId == r.Id {
				(*rooms)[i].Messages = append((*rooms)[i].Messages, internal.Message{
					Id:           uuid.New().String(),
					Content:      data.Content,
					RoomId:       data.RoomId,
					ConnectionId: conn.Id,
					DeliveredTo:  make(map[string]bool),
				})
			}

			log.Printf("New message received from: %s. Content: %s", conn.Name, data.Content)
		}
	}
}

func socket_conn(w http.ResponseWriter, r *http.Request, rooms *[]internal.Room) {
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

	data := unmarshal_message[internal.Connection](c, mt, buf)
	if data == nil {
		return
	}

	connId := uuid.New().String()
	conn := internal.Connection{
		Id:     connId,
		Name:   data.Name,
		RoomId: data.RoomId,
	}

	for i, r := range *rooms {
		if r.Id == conn.RoomId {
			(*rooms)[i].Connections = append((*rooms)[i].Connections, conn)
		}
	}

	res, err := json.Marshal(data)
	if err != nil {
		log.Println("error: ", err)
		return
	}
	log.Println(res)
	c.WriteMessage(mt, res)
	log.Printf("Connected: %s", conn.Name)

	go read_messages(c, conn, rooms)
	go broadcast_messages(c, conn, rooms)
}

const port = "8080"

func main() {
	flag.Parse()
	var rooms []internal.Room
	const mockedRoomId string = "21ca15d0-e346-4630-a240-773a828c31b3"
	rooms = append(rooms, internal.Room{
		Id: mockedRoomId,
	})

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) { socket_conn(w, r, &rooms) })
	log.Printf("Server running at port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), nil))
}
