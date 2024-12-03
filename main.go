package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/YzmjY/qq-bot/internal/agent"
)

var (
	ga agent.Agent
)

func EventHandler(ctx context.Context, event *agent.Event) error {
	if event.MessageType != "group" {
		return nil
	}
	if event.UserID == int(event.SelfID) {
		return nil
	}

	sender := event.Sender.UserID

	return ga.SendGroupMsg(ctx, int64(event.GroupID), []byte("接收了来自"+fmt.Sprintf("%d", sender)+"的消息"))
}

func main() {
	a, err := agent.NewNapcatAgent(context.Background(), map[string]string{
		"server_addr": "127.0.0.1",
		"server_port": "8080",
		"event_addr":  "127.0.0.1",
		"event_port":  "3000",
	})
	if err != nil {
		panic(err)
	}

	ga = a

	serveForEvent()
}
func serveForEvent() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.Buffer{}
		io.Copy(&buf, r.Body)
		fmt.Printf("receive event: %s\n", buf.String())
		var e agent.Event
		_ = json.Unmarshal(buf.Bytes(), &e)
		EventHandler(context.Background(), &e)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", 8080), handler)
}
