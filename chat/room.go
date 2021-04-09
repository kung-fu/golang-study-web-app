package main

type room struct {
	// 他のクライアントに転送するメッセージを保持するチャネル
	forward chan []byte
	// チャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	// チャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	// 在室しているすべてのクライアントを保持
	clients map[*client]bool
}
