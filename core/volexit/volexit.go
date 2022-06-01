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

package volexit

import (
	"context"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/leolara/go-event"
)

type ParSignedVoluntaryExit struct {
	Exit     phase0.SignedVoluntaryExit
	ShareIdx uint
}

type ParSigRepo interface {
	StoreInternal(ctx context.Context, parSig ParSignedVoluntaryExit) error
	SubscribeThreshold(channel chan<- []ParSignedVoluntaryExit) event.Subscription
	Shutdown()
}

type ParSigExchange interface {
	Broadcast(ctx context.Context, parSig ParSignedVoluntaryExit)
	Subscribe(chan<- ParSignedVoluntaryExit) event.Subscription
}

type voluntaryExitParSigExchange struct {
}

func (voluntaryExitParSigExchange) Broadcast(ctx context.Context, parSig ParSignedVoluntaryExit) {

}

func (voluntaryExitParSigExchange) Subscribe() (chan<- []ParSignedVoluntaryExit, error) {
	return nil, nil
}

type voluntaryExitSigAggregator struct {
	parSig ParSigRepo
}

func (voluntaryExitSigAggregator) aggregate([]ParSignedVoluntaryExit) error {
	return nil
}

func (voluntaryExitSigAggregator) Subscribe() (chan<- phase0.SignedVoluntaryExit, error) {
	return nil, nil
}

type voluntaryExitBroadcaster struct {
	agg *voluntaryExitSigAggregator
}

func (voluntaryExitBroadcaster) BroadcastVoluntaryExit(phase0.SignedVoluntaryExit) error {
	return nil
}
