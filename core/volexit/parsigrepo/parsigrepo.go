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

package parsigrepo

import (
	"context"
	"fmt"
	"sync"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/leolara/go-event"

	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/core/volexit"
)

type parSigRepo struct {
	exchange  volexit.ParSigExchange
	threshold uint

	entries       map[phase0.ValidatorIndex][]volexit.ParSignedVoluntaryExit
	mutex         sync.Mutex
	exchangeCh    chan volexit.ParSignedVoluntaryExit
	exchangeSub   event.Subscription
	shutdown      chan struct{}
	thresholdFeed event.Feed[[]volexit.ParSignedVoluntaryExit]
}

func New(
	exchange volexit.ParSigExchange,
	threshold uint,
) volexit.ParSigRepo {
	repo := &parSigRepo{
		exchange:  exchange,
		threshold: threshold,
	}
	repo.init()

	return repo
}

func (repo *parSigRepo) StoreInternal(ctx context.Context, parSig volexit.ParSignedVoluntaryExit) error {
	err := repo.store(parSig)
	if err != nil {
		return errors.Wrap(err, "volexit/parsigrepo: storing parsig")
	}

	repo.exchange.Broadcast(ctx, parSig)

	return nil
}

func (repo *parSigRepo) SubscribeThreshold(channel chan<- []volexit.ParSignedVoluntaryExit) event.Subscription {
	return repo.thresholdFeed.Subscribe(channel)
}

func (repo *parSigRepo) Shutdown() {
	close(repo.shutdown)
}

func (repo *parSigRepo) init() {
	repo.entries = make(map[phase0.ValidatorIndex][]volexit.ParSignedVoluntaryExit)
	repo.exchangeCh = make(chan volexit.ParSignedVoluntaryExit, 10)
	repo.exchangeSub = repo.exchange.Subscribe(repo.exchangeCh)
	repo.shutdown = make(chan struct{})

	go func() {
		for {
			select {
			case ve := <-repo.exchangeCh:
				err := repo.StoreInternal(context.Background(), ve)
				if err != nil {
					// TODO log error
				}
			case err := <-repo.exchangeSub.Err():
				if err != nil {
					// TODO log error
				}
			case <-repo.shutdown:
				repo.exchangeSub.Unsubscribe()
				return
			}
		}
	}()
}

func (repo *parSigRepo) storeExternal(parSig volexit.ParSignedVoluntaryExit) error {
	err := repo.store(parSig)
	if err != nil {
		return errors.Wrap(err, "volexit/parsigrepo: storing parsig")
	}

	return nil
}

func (repo *parSigRepo) store(parSig volexit.ParSignedVoluntaryExit) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	valIdx := parSig.Exit.Message.ValidatorIndex

	existing, ok := repo.entries[valIdx]
	if ok && len(existing) > 0 {
		curEpoch := existing[0].Exit.Message.Epoch
		newEpoch := parSig.Exit.Message.Epoch
		if curEpoch != newEpoch {
			return fmt.Errorf(
				"epoch mismatch in parsig of voluntary exit for validator index %d (%d != %d)",
				valIdx,
				curEpoch,
				newEpoch,
			)
		}
	}

	if repo.existsShareIdx(valIdx, parSig.ShareIdx) {
		// idempotent if inserted twice
		return nil
	}

	repo.entries[valIdx] = append(existing, parSig)

	return nil
}

func (repo *parSigRepo) existsShareIdx(valIdx phase0.ValidatorIndex, shareIdx uint) bool {
	existing, ok := repo.entries[valIdx]
	if !ok || len(existing) == 0 {
		return false
	}

	for _, entry := range existing {
		if entry.ShareIdx == shareIdx {
			return true
		}
	}

	return false
}

func (repo *parSigRepo) isThresholdReached(valIdx phase0.ValidatorIndex) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	return len(repo.entries[valIdx]) == int(repo.threshold)
}
