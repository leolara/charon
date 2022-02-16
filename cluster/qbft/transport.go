package qbft

import (
	"context"
	"github.com/corverroos/quorum/consensus/istanbul"
	"github.com/corverroos/quorum/event"
	"github.com/obolnetwork/charon/app/log"
	"sync"
)

type Transport interface {
	// Subscribe subscribes to messages as index.
	Subscribe(index int, mux *event.TypeMux)
	// SendAll sends a message to all validators (exclude index).
	SendAll(index int, msg istanbul.MessageEvent) error
}

func NewMemTransport() Transport {
	return &memTransport{muxs: make(map[int]*event.TypeMux)}
}

type memTransport struct {
	mu   sync.Mutex
	muxs map[int]*event.TypeMux
}

func (m *memTransport) SendAll(index int, msg istanbul.MessageEvent) error {
	go func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		for i, mux := range m.muxs {
			if i == index {
				continue
			}

			err := mux.Post(msg)
			if err != nil {
				log.Error(context.TODO(), "transport post error", err)
			}
		}
	}()

	return nil
}

func (m *memTransport) Subscribe(index int, mux *event.TypeMux) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.muxs[index] = mux
}
