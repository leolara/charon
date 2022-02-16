package qbft

import (
	"github.com/corverroos/quorum/common"
	"github.com/corverroos/quorum/consensus/istanbul"
	qbfttypes "github.com/corverroos/quorum/consensus/istanbul/qbft/types"
	"github.com/corverroos/quorum/rlp"
	"github.com/obolnetwork/charon/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestProposal(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	payload := make([]byte, 8)
	random.Read(payload)

	p1 := proposal{
		Duty: types.Duty{
			Slot: rand.Int63(),
			Type: types.DutyType(rand.Int()),
		},
		Data: payload,
	}
	raw1 := p1.Raw()

	p2, err := proposalFromRaw(raw1)
	require.NoError(t, err)

	raw2 := p2.Raw()

	require.Equal(t, p1, p2)
	require.Equal(t, raw1, raw2)
}

func TestProposalRLP(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	payload := make([]byte, 8)
	random.Read(payload)

	p1 := newProposal(random)

	b, err := rlp.EncodeToBytes(p1)
	require.NoError(t, err)

	var p2 istanbul.RawProposal
	err = rlp.DecodeBytes(b, &p2)
	require.NoError(t, err)

	require.Equal(t, p1.Number().Uint64(), p2.Number().Uint64())
	require.Equal(t, p1.Hash(), p2.Hash())
	require.Equal(t, p1.String(), p2.String())
	require.Equal(t, p1.Payload(), p2.Payload())
}

func TestPreprepareRLP(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	payload := make([]byte, 8)
	random.Read(payload)

	pp1 := qbfttypes.NewPreprepare(
		big.NewInt(random.Int63()),
		big.NewInt(random.Int63()),
		newProposal(random),
	)

	b, err := rlp.EncodeToBytes(pp1)
	require.NoError(t, err)

	msg, err := qbfttypes.Decode(pp1.Code(), b)
	require.NoError(t, err)

	pp2, ok := msg.(*qbfttypes.Preprepare)
	require.True(t, ok)

	require.Equal(t, pp1.Proposal, pp2.Proposal)
	require.Equal(t, pp1.Round, pp2.Round)
	require.Equal(t, pp1.Sequence, pp2.Sequence)
	require.Equal(t, pp1, pp2)
}

func newProposal(random *rand.Rand) *istanbul.RawProposal {
	payload := make([]byte, 8)
	random.Read(payload)

	return istanbul.NewRawProposal(
		big.NewInt(random.Int63()),
		common.BigToHash(big.NewInt(random.Int63())),
		payload,
		"label")
}
