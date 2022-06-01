// Copyright © 2022 Obol Labs Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along with
// this program.  If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"context"
	"encoding/json"

	"github.com/leolara/go-event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/app/log"
	"github.com/obolnetwork/charon/app/z"
)

func NewP2PTransport[T any](protocol protocol.ID, tcpNode host.Host, peerIDx int, peers []peer.ID) interface{} {
	t := &genericP2PTransport[T]{
		protocol: protocol,
		tcpNode:  tcpNode,
		peers:    peers,
		peerIDx:  peerIDx,
		ch:       make(chan p2pMsg[T]),
	}
	t.tcpNode.SetStreamHandler(protocol, t.handle)

	return t
}

type genericP2PTransport[T any] struct {
	tcpNode  host.Host
	peerIDx  int
	peers    []peer.ID
	ch       chan p2pMsg[T]
	protocol protocol.ID

	incomming event.Feed[p2pMsg[T]]
}

// handle implements p2p network.StreamHandler processing new incoming messages.
func (t *genericP2PTransport[T]) handle(s network.Stream) {
	defer s.Close()

	var msg p2pMsg[T]
	err := json.NewDecoder(s).Decode(&msg)
	if err != nil {
		log.Error(context.Background(), "Decode leadercast message", err)
		return
	}
	t.incomming.Send(msg)
}

func (t *genericP2PTransport[T]) Broadcast(ctx context.Context, fromIdx int, data T) error {
	b, err := json.Marshal(p2pMsg[T]{
		FromIdx: fromIdx,
		Data:    data,
	})
	if err != nil {
		return errors.Wrap(err, "marshal tcpNode msg")
	}

	var errs []error

	for idx, p := range t.peers {
		if idx == t.peerIDx {
			// Don't send to self.
			continue
		}

		err := t.sendData(ctx, p, b)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		log.Debug(ctx, "Leader broadcast duty success")
	} else {
		// len(t.peers)-len(errs)-1 is total number of errors excluding broadcast to self case
		log.Warn(ctx, "Broadcast duty with errors", errs[0], z.Int("success", len(t.peers)-len(errs)-1),
			z.Int("errors", len(errs)))
	}

	return nil
}

func (t *genericP2PTransport[T]) Subscribe(channel chan<- p2pMsg[T]) event.Subscription {
	return t.incomming.Subscribe(channel)
}

func (t *genericP2PTransport[T]) sendData(ctx context.Context, p peer.ID, b []byte) error {
	// Circuit relay connections are transient
	s, err := t.tcpNode.NewStream(network.WithUseTransient(ctx, "leadercast"), p, t.protocol)
	if err != nil {
		return errors.Wrap(err, "tcpNode stream")
	}

	_, err = s.Write(b)
	if err != nil {
		return errors.Wrap(err, "tcpNode write")
	}

	if err := s.Close(); err != nil {
		return errors.Wrap(err, "tcpNode close")
	}

	return nil
}

type p2pMsg[T any] struct {
	Idx     int
	FromIdx int
	Data    T
}
