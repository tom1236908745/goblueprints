package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// clientはチャットを行っているある１人のユーザーを表している。
type client struct {
	// socketはクライアント用
	socket *websocket.Conn
	// sendはメッセージが遅れるチャンネル
	send chan *message
	// roomはこのクライアントが参加しているチャットルーム
	room *room
	// useDataはユーザーの情報を持っている
	userData map[string]interface{}
}

// クライアントがWebSocketに読み書きするためのメソッド
func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now().Format("2006-01-02 03:04:05")
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	// 継続的にsendチャンネルからのメッセージを受け取る。
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
