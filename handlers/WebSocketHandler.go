package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

// global channel untuk komunikasi websocket payload dan goroutine
var wsChan = make(chan WsPayload)

// map yang digunakan untuk mennyimpan koneksi yang dapat digunakan untuk indikasi user
var clients = make(map[WebSocketConnection]string)

// spesifikasi upgrade koneksi ke socket (full suplex)
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// format json response untuk koneksi web socket (response dari server)
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// format payload untuk pesan di websocket (Request dari client)
type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// struct koneksi web socket
type WebSocketConnection struct {
	*websocket.Conn
}

// function / endpoint untuk mengupgrade koneksi
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// menjalankan upgrade koneksi
	ws, err := upgradeConnection.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	// mengirim response
	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small/></em>`

	//mencatat koneksi yang masuk
	conn := WebSocketConnection{ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

// function untuk listen pesan websocket pada koneksi tertentu
func ListenForWs(conn *WebSocketConnection) {
	// menjalankan function ini ketika selesai eksekusi function listenforws, untuk memastikan function benar benar berhenti
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {

		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// function untuk melisten global channel untuk pesan yang akan datang
func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_user"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			response.Action = "list_user"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcastToAll(response)
		}

		//response.Action = "Got here"
		//response.Message = fmt.Sprintf("Action %s", e.Action)
	}
}

// mendistribusikan pesan ke setiap user yang terkoneksi
func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}

	sort.Strings(userList)
	return userList
}
