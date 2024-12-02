package agent

import (
	"context"
	"errors"
)

var (
	AgentFactories = make(map[string]AgentFactory)

	ErrNotSupportAgentType = errors.New("not support agent type")
)

type Agent interface {
}

type EventCb func(ctx context.Context, eventRaw []byte) error

type AgentFactory interface {
	CreateAgent(ctx context.Context, params map[string]string) (Agent, error)
}

func RegisterAgent(typ string, factory AgentFactory) {
	AgentFactories[typ] = factory
}

func NewAgent(ctx context.Context, typ string, params map[string]string) (Agent, error) {
	factory, ok := AgentFactories[typ]
	if !ok {
		return nil, ErrNotSupportAgentType
	}

	return factory.CreateAgent(ctx, params)
}
