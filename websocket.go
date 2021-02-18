package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("[SOCKET] - Client Connected")

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					log.Printf("[SOCKET] - Client Disconnected: %s", err)
				}
			}
			break
		}
		if len(p) < 2 || !strings.HasPrefix(string(p), "{") || !strings.HasSuffix(string(p), "}") {
			continue
		}

		var v = wsRequestPayload{}

		if err := json.Unmarshal(p, &v); err != nil {
			log.Printf("[ERROR: Decode JSON] - %q", err)
			log.Println("[SOCKET] - Server disconnected client")
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1003, "Invalid JSON Received"))
			break
		}

		if len(v.Method) == 0 || !isValidMethod(v.Method) {
			sendMessage(conn, messageType, "{\"type\": \"warn\", \"message\": \"Warning: valid method property must be present in JSON payloads\"}")
		}

		if v.Method != "SET" && v.Expires != 0 {
			sendMessage(conn, messageType, "{\"type\": \"warn\", \"message\": \"Warning: expires property should not be present when method is "+v.Method+"\"}")
			v.Expires = 0
		}
		if v.Method != "SET" && len(v.Value) != 0 {
			sendMessage(conn, messageType, "{\"type\": \"warn\", \"message\": \"Warning: value property should not be present when method is "+v.Method+"\"}")
			v.Value = ""
		}
		if v.Method != "CLEAR" && len(v.Key) == 0 {
			sendMessage(conn, messageType, "{\"type\": \"error\", \"message\": \"Warning: key property must be present when method is "+v.Method+"\"}")
			continue
		}

		if v.Method == "SET" {
			if len(v.Value) == 0 {
				sendMessage(conn, messageType, "{\"type\": \"error\", \"message\": \"Warning: value property must be present when method is "+v.Method+"\"}")
				continue
			}

			if v.Expires > 0 && v.Expires < time.Now().Unix() * 1000 {
				sendMessage(conn, messageType, "{\"type\": \"warn\", \"message\": \"Warning: expires property should be greater than current time\"}")
				v.Expires = 0
			}

			setCache(v.Key, v.Value, v.Expires)
		} else if v.Method == "DELETE" {
			deleteCache(v.Key)
		} else if v.Method == "GET" {
			val := getCache(v.Key)

			if len(val) == 0 {
				val = "null"
			} else {
				val = `"` + val + `"`
			}

			sendMessage(conn, messageType, "{\"type\": \"response\", \"key\": \""+v.Key+"\", \"value\":"+val+"}")
		} else if v.Method == "CLEAR" {
			clearCache()
		}
	}
}

func sendMessage(c *websocket.Conn, m int, s string) {
	if err := c.WriteMessage(m, []byte(s)); err != nil {
		log.Println(err)
	}
}
