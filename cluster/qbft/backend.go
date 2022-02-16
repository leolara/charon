package qbft

import (
	"context"
	"crypto/ecdsa"
	"github.com/corverroos/quorum/common"
	"github.com/corverroos/quorum/consensus/istanbul"
	"github.com/corverroos/quorum/consensus/istanbul/qbft/core"
	"github.com/corverroos/quorum/consensus/istanbul/validator"
	"github.com/corverroos/quorum/event"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obolnetwork/charon/app/log"
	"github.com/obolnetwork/charon/app/z"
	"github.com/obolnetwork/charon/types"
	"math/big"
	"time"
)

func New(transport Transport, index int, pubkeys []ecdsa.PublicKey, privkey *ecdsa.PrivateKey) backend {
	policy := istanbul.NewRoundRobinProposerPolicy()

	b := backend{
		index:     index,
		pubkeys:   pubkeys,
		privkey:   privkey,
		dutyType:  types.DutyAttester,
		transport: transport,
		policy:    policy,
		eventMux:  new(event.TypeMux),
		commmitFunc: func(proposal istanbul.Proposal, seals [][]byte, round *big.Int) error {
			log.Info(context.TODO(), "commit", z.Any("proposal", proposal), z.U64("round", round.Uint64()))
			return nil
		},
	}

	config := istanbul.DefaultConfig
	config.ProposerPolicy = policy

	b.core = core.New(&b, config)

	transport.Subscribe(index, b.EventMux())

	go func() {
		time.Sleep(time.Second)

		err := b.EventMux().Post(istanbul.RequestEvent{
			Proposal: proposal{
				Duty: types.Duty{Slot: int64(123), Type: types.DutyAttester},
				Data: []byte{byte(index)},
			}.Raw(),
		})
		if err != nil {
			panic(err)
		}

	}()

	return b
}

type backend struct {
	index       int
	pubkeys     []ecdsa.PublicKey
	dutyType    types.DutyType
	eventMux    *event.TypeMux
	transport   Transport
	policy      *istanbul.ProposerPolicy
	commmitFunc func(proposal istanbul.Proposal, seals [][]byte, round *big.Int) error
	privkey     *ecdsa.PrivateKey
	core        istanbul.Core
}

func (b backend) Start() error {
	return b.core.Start()
}

func (b backend) Stop() error {
	return b.core.Stop()
}

// Address returns the owner's address.
// It returns the charon node's cluster index.
func (b backend) Address() common.Address {
	return common.Address(crypto.PubkeyToAddress(b.privkey.PublicKey))
}

// Validators returns the validator set.
// It returns the charon cluster.
func (b backend) Validators(istanbul.Proposal) istanbul.ValidatorSet {
	var addrs []common.Address
	for _, pubkey := range b.pubkeys {
		addrs = append(addrs, common.Address(crypto.PubkeyToAddress(pubkey)))
	}
	return validator.NewSet(addrs, b.policy)
}

// EventMux returns the event mux in backend.
func (b backend) EventMux() *event.TypeMux {
	return b.eventMux
}

// Broadcast sends a message to all validators (include self).
func (b backend) Broadcast(_ istanbul.ValidatorSet, code uint64, payload []byte) error {
	msg := istanbul.MessageEvent{
		Code:    code,
		Payload: payload,
	}

	err := b.transport.SendAll(-1, msg)
	if err != nil {
		return err
	}

	return nil
}

// Gossip sends a message to all validators (exclude self)
func (b backend) Gossip(_ istanbul.ValidatorSet, code uint64, payload []byte) error {
	msg := istanbul.MessageEvent{
		Code:    code,
		Payload: payload,
	}

	err := b.transport.SendAll(b.index, msg)
	if err != nil {
		return err
	}

	return nil
}

// Commit delivers an approved proposal to backend.
// The delivered proposal will be put into blockchain.
func (b backend) Commit(proposal istanbul.Proposal, seals [][]byte, round *big.Int) error {
	return b.commmitFunc(proposal, seals, round)
}

// Verify verifies the proposal. If a consensus.ErrFutureBlock error is returned,
// the time difference of the proposal and current time is also returned.
func (b backend) Verify(proposal istanbul.Proposal) (time.Duration, error) {
	//TODO implement me
	return 0, nil
}

// Sign signs input data with the backend's private privkey.
func (b backend) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(data), b.privkey)
}

// SignWithoutHashing sign input data with the backend's private key without hashing the input data
func (b backend) SignWithoutHashing(digest []byte) ([]byte, error) {
	return crypto.Sign(digest, b.privkey)
}

// CheckSignature verifies the signature by checking if it's signed by
// the given validator
func (b backend) CheckSignature(data []byte, addr common.Address, sig []byte) error {
	//TODO implement me
	return nil
}

// LastProposal retrieves latest committed proposal and the address of proposer
func (b backend) LastProposal() (istanbul.Proposal, common.Address) {
	return istanbul.NewRawProposal(
		big.NewInt(0),
		common.HexToHash(""),
		nil,
		"",
	), common.BytesToAddress(nil)
}

// HasPropsal checks if the combination of the given hash and height matches any existing blocks
func (b backend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return false
}

// GetProposer returns the proposer of the given block height
func (b backend) GetProposer(number uint64) common.Address {
	a := common.Address(crypto.PubkeyToAddress(b.pubkeys[int(number)%len(b.pubkeys)]))
	return a
}

// ParentValidators returns the validator set of the given proposal's parent block
func (b backend) ParentValidators(proposal istanbul.Proposal) istanbul.ValidatorSet {
	return b.Validators(proposal)
}

func (b backend) HasBadProposal(hash common.Hash) bool {
	return false
}

func (b backend) Close() error {
	return nil
}

func (b backend) IsQBFTConsensusAt(*big.Int) bool {
	return true
}

func (b backend) StartQBFTConsensus() error {
	return nil
}
