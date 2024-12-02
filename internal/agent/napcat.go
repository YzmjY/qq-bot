package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const Name = "napcat"

type napcatConfig struct {
	// 主动向服务端发送消息的地址
	ServerAddr  string `yaml:"server_addr"`
	ServerPort  string `yaml:"server_port"`
	ServerToken string `yaml:"server_token"`

	// 接收服务端推送事件的地址
	EventAddr  string `yaml:"event_addr"`
	EventPort  string `yaml:"event_port"`
	EventToken string `yaml:"event_token"`
}

type NapcatAgent struct {
	napcatConfig

	eventCB EventCb
	eventCh chan []byte

	ctx    context.Context
	cancel context.CancelFunc
}

func NewNapcatAgent(ctx context.Context, params map[string]string) (Agent, error) {
	return &NapcatAgent{
		napcatConfig: napcatConfig{
			ServerAddr:  params["server_addr"],
			ServerPort:  params["server_port"],
			ServerToken: params["server_token"],
			EventAddr:   params["event_addr"],
			EventPort:   params["event_port"],
			EventToken:  params["event_token"],
		},
	}, nil
}

func (a *NapcatAgent) Start(ctx context.Context) error {
	return nil
}

func (a *NapcatAgent) activate() {
	a.serveForEvent()
}

func (a *NapcatAgent) serveForEvent() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.Buffer{}
		io.Copy(&buf, r.Body)

		a.eventCh <- buf.Bytes()
	})

	http.ListenAndServe(fmt.Sprintf("%s:%s", a.EventAddr, a.EventPort), handler)
}

func (a *NapcatAgent) SendMsg()
