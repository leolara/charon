package qbft

import (
	"encoding/json"
	"fmt"
	"github.com/corverroos/quorum/common"
	"github.com/corverroos/quorum/consensus/istanbul"
	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/types"
	"golang.org/x/crypto/sha3"
	"math/big"
)

func proposalFromRaw(p istanbul.Proposal) (proposal, error) {
	raw, ok := p.(*istanbul.RawProposal)
	if !ok {
		return proposal{}, errors.New("invalid proposal type")
	}

	var resp proposal
	err := json.Unmarshal(raw.Payload(), &resp)
	if err != nil {
		return proposal{}, err
	}

	return resp, nil
}

type proposal struct {
	Duty types.Duty
	Data []byte
}

// Raw returns the proposal in raw format.
func (p proposal) Raw() *istanbul.RawProposal {
	payload, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return istanbul.NewRawProposal(
		big.NewInt(1), // Always sequence/height 1
		p.Hash(),
		payload,
		p.String(),
	)
}

// Hash retrieves the hash of this proposal.
func (p proposal) Hash() common.Hash {
	hasher := sha3.New256()

	if _, err := hasher.Write(big.NewInt(p.Duty.Slot).Bytes()); err != nil {
		panic(err)
	}

	if _, err := hasher.Write(big.NewInt(int64(p.Duty.Type)).Bytes()); err != nil {
		panic(err)
	}

	if _, err := hasher.Write(p.Data); err != nil {
		panic(err)
	}

	return common.BytesToHash(hasher.Sum(nil))
}

func (p proposal) String() string {
	return fmt.Sprintf("%s/%4x", p.Duty, p.Data)
}
