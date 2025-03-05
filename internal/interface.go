package internal

type Connection struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	RoomId string `json:"room_id"`
}

type Message struct {
	Id             string          `json:"id"`
	Content        string          `json:"content"`
	DeliveredTo    map[string]bool `json:"delivered_to"`
	DeliveredCount int             `json:"delivered_count"`
	ConnectionId   string          `json:"connection_id"`
	RoomId         string          `json:"room_id"`
}

type Room struct {
	Id          string       `json:"id"`
	Connections []Connection `json:"connections"`
	Messages    []Message    `json:"messages"`
}
