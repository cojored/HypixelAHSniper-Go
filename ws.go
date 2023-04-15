package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

var sockets = make(map[string]*websocket.Conn)

func generateId(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func socket(ws *websocket.Conn) {
	id := generateId(15)
	sockets[id] = ws

	/*go func() {
		for {
			var msg = make([]byte, 512)
			_, _ = ws.Read(msg)
			fmt.Println(string(msg[:]))
		}
	}()*/

	ticker := time.NewTicker(3 * time.Second)
	func() {
		for {
			select {
			case <-ticker.C:
				_, err := ws.Write([]byte("ping"))
				if err != nil {
					delete(sockets, id)
				}
			}
		}
	}()
}

func sendMessage(message string) {
	for id, ws := range sockets {
		_, err := ws.Write([]byte(message))
		if err != nil {
			delete(sockets, id)
		}
	}
}
