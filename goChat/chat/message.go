package main

// messageでメッセージに持たせたい情報を格納する。
type message struct {
	Name      string // ユーザー名
	Message   string // メッセージ
	When      string // 送信した時刻
	AvatarURL string // 画像
}
