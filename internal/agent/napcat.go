package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type Agent interface {
	HandleEvent(ctx context.Context, event *Event) error
	SendGroupMsg(ctx context.Context, groupID int64, msg []byte) error
}

type Event struct {
	SelfID      int64  `json:"self_id"`
	UserID      int    `json:"user_id"`
	Time        int    `json:"time"`
	MessageID   int    `json:"message_id"`
	MessageSeq  int    `json:"message_seq"`
	RealID      int    `json:"real_id"`
	MessageType string `json:"message_type"`
	Sender      struct {
		UserID   int    `json:"user_id"`
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
		Role     string `json:"role"`
	} `json:"sender"`
	RawMessage string `json:"raw_message"`
	Font       int    `json:"font"`
	SubType    string `json:"sub_type"`
	Message    []struct {
		Type string `json:"type"`
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
	} `json:"message"`
	MessageFormat string `json:"message_format"`
	PostType      string `json:"post_type"`
	GroupID       int    `json:"group_id"`
}

type Handler func(ctx context.Context, event *Event) error

const Name = "napcat"

type napcatConfig struct {
	// 主动向服务端发送消息的地址
	ServerAddr  string `yaml:"server_addr"`
	ServerPort  string `yaml:"server_port"`
	ServerToken string `yaml:"server_token"`
}

type NapcatAgent struct {
	napcatConfig

	eventHandler []Handler
	eventCh      chan []byte

	ctx    context.Context
	cancel context.CancelFunc

	cli *http.Client
}

func NewNapcatAgent(ctx context.Context, params map[string]string) (Agent, error) {
	agent := &NapcatAgent{
		napcatConfig: napcatConfig{
			ServerAddr: params["server_addr"],
			ServerPort: params["server_port"],
		},
	}

	agent.ctx, agent.cancel = context.WithCancel(ctx)
	agent.eventCh = make(chan []byte, 1024)

	return agent, nil
}

func (a *NapcatAgent) HandleEvent(ctx context.Context, event *Event) error {
	for _, h := range a.eventHandler {
		if err := h(ctx, event); err != nil {
			log.Error().Err(err).Msg("handle event error")
		}
	}
	return nil
}

func (a *NapcatAgent) SendMsg(ctx context.Context, msg []byte) error {
	return nil
}

type TextMessage struct {
	Type string      `json:"type"`
	Data TextMsgData `json:"data"`
}

type TextMsgData struct {
	Text string `json:"text"`
}
type GroupMsg struct {
	GroupID int64         `json:"group_id"`
	Message []TextMessage `json:"message"`
}

func (a *NapcatAgent) SendGroupMsg(ctx context.Context, groupID int64, msg []byte) error {
	url := "http://localhost:3000/send_group_msg"
	method := "POST"

	groupMsg := GroupMsg{
		GroupID: groupID,
		Message: []TextMessage{
			{
				Type: "text",
				Data: TextMsgData{
					Text: string(msg),
				},
			},
		},
	}
	all, err := json.Marshal(groupMsg)
	if err != nil {
		log.Error().Err(err).Msg("marshal group msg error")
	}

	payload := strings.NewReader(string(all))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("send group msg error")
		return err
	}
	defer res.Body.Close()

	return nil
}
