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

		if len(v.Op) == 0 || !isValidOp(v.Op) {
			sendMessage(conn, messageType, "{\"op\": \"WARN\", \"message\": \"Warning: valid op property must be present in JSON payloads\"}")
		}

		if v.Op != "SET" && v.Expires != 0 {
			sendMessage(conn, messageType, "{\"op\": \"WARN\", \"message\": \"Warning: expires property should not be present when op is "+v.Op+"\"}")
			v.Expires = 0
		}
		if v.Op != "SET" && len(v.Value) != 0 {
			sendMessage(conn, messageType, "{\"op\": \"WARN\", \"message\": \"Warning: value property should not be present when op is "+v.Op+"\"}")
			v.Value = ""
		}
		if v.Op != "CLEAR" && len(v.Key) == 0 {
			sendMessage(conn, messageType, "{\"op\": \"ERROR\", \"message\": \"Warning: key property must be present when op is "+v.Op+"\"}")
			continue
		}

		if v.Op == "SET" {
			if len(v.Value) == 0 {
				sendMessage(conn, messageType, "{\"op\": \"ERROR\", \"message\": \"Warning: value property must be present when op is "+v.Op+"\"}")
				continue
			}

			if v.Expires > 0 && v.Expires < time.Now().Unix()*1000 {
				sendMessage(conn, messageType, "{\"op\": \"WARN\", \"message\": \"Warning: expires property should be greater than current time\"}")
				v.Expires = 0
			}

			setCache(v.Key, v.Value, v.Expires)
		} else if v.Op == "DELETE" {
			deleteCache(v.Key)
		} else if v.Op == "GET" {
			val := getCache(v.Key)

			if len(val) == 0 {
				val = "null"
			} else {
				val = `"` + val + `"`
			}

			sendMessage(conn, messageType, "{\"op\": \"RESPONSE\", \"key\": \""+v.Key+"\", \"value\":"+val+"}")
		} else if v.Op == "CLEAR" {
			clearCache()
		}
	}
}

func sendMessage(c *websocket.Conn, m int, s string) {
	if err := c.WriteMessage(m, []byte(s)); err != nil {
		log.Println(err)
	}
}
