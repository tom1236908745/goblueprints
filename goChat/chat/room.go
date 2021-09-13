package main

import (
	"log"
	"net/http"

	"github.com/stretchr/objx"

	"github.com/gorilla/websocket"
	"github.com/tom1236908745/goChat/trace"
)

type room struct {
	// forwardは他のクライアントに転送するためのメッセージを保持するためのチャンネル
	forward chan *message
	//　joinはチャットルームに参加しようとしているクライアントのためのチャンネル
	join chan *client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャンネル
	leave chan *client
	// clientsには在室している全てのクライアントが保持されている。
	clients map[*client]bool
	// tracerはチャットルーム上で行われた操3作のログを取る。
	tracer trace.Tracer
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

// runメソッド内で、join, leave, clientsを監視 (selectのcase節ではこれらが同時に実行される事はない)
func (r *room) run() {
	for {
		select {
		// 参加
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました。")
			// 退室
		case client := <-r.leave:
			delete(r.clients, client)
			// チャンネルを閉じる
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", msg.Message)
			// 全てのクライアントにメッセージを送信
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
					r.tracer.Trace(" -- クライアントに送信されました")
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします。")
				}
			}
		}
	}
}

const (
	socketBufferSize   = 1024
	messageBuffferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました。:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBuffferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
