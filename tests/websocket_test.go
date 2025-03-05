package client

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/WilliamKSilva/type-1v1/internal"
	"github.com/gorilla/websocket"
)

const mockedRoomId string = "21ca15d0-e346-4630-a240-773a828c31b3"

func TestWebsocketConnection(t *testing.T) {
	c, err := internal.EstablishConnection("Teste", mockedRoomId)
	if err != nil {
		t.Fatal("Error trying to estabilish connection:", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Log("read: ", err)
			return
		}

		var connection internal.Connection
		err = json.Unmarshal(message, &connection)
		if err != nil {
			return
		}
		if connection.Name != "" {
			t.Log("Connection established:", connection.Name)
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}

// TODO: Servidor crashando quando se roda esse teste 2 vezes seguidas
func TestWebsocketSendMessageToOtherConnection(t *testing.T) {
	cSender, err := internal.EstablishConnection("CSender", mockedRoomId)
	if err != nil {
		t.Fatal("Error trying to estabilish connection cSender:", err)
		return
	}

	cReceiver, err := internal.EstablishConnection("CReceiver", mockedRoomId)
	if err != nil {
		t.Fatal("Error trying to estabilish connection cReceiver:", err)
		return
	}

	defer cSender.Close()
	defer cReceiver.Close()

	mes := internal.Message{
		Content: "Sending message test",
		RoomId:  mockedRoomId,
	}
	data, err := json.Marshal(mes)
	if err != nil {
		t.Fatal("Websocket error trying to marshal Message:", err)
		return
	}

	err = cSender.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		t.Fatal("Websocket error trying to send Message data:", err)
		return
	}

	for {
		_, m, err := cReceiver.ReadMessage()
		if err != nil {
			t.Log("read: ", err)
			return
		}

		var message internal.Message
		json.Unmarshal(m, &message)
		log.Println("Received message content:", message.Content)
	}
}
