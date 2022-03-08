// Copyright © 2021 Obol Technologies Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sigagg_test

import (
	"context"
	"crypto/rand"
	"testing"

	bls12381 "github.com/dB2510/kryptology/pkg/core/curves/native/bls12-381"
	"github.com/dB2510/kryptology/pkg/signatures/bls/bls_sig"
	"github.com/stretchr/testify/require"

	"github.com/obolnetwork/charon/core"
	"github.com/obolnetwork/charon/core/sigagg"
	"github.com/obolnetwork/charon/tbls"
	"github.com/obolnetwork/charon/testutil"
)

func TestSigAgg(t *testing.T) {
	ctx := context.Background()

	const (
		threshold = 3
		peers     = 4
	)

	att := testutil.RandomAttestation()

	// Sign the attestation directly (spec domain not required for test)
	msg, err := att.MarshalSSZ()
	require.NoError(t, err)

	// Generate private shares
	tss, secrets, err := tbls.GenerateTSS(threshold, peers, rand.Reader)
	require.NoError(t, err)

	// Create partial signatures (in two formats)
	var (
		parsigs []core.ParSignedData
		psigs   []*bls_sig.PartialSignature
	)
	for _, secret := range secrets {
		sig, err := tbls.PartialSign(secret, msg)
		require.NoError(t, err)

		copy(att.Signature[:], bls12381.NewG2().ToCompressed(sig.Signature))

		parsig, err := core.EncodeAttestationParSignedData(att, int(sig.Identifier))
		require.NoError(t, err)

		psigs = append(psigs, sig)
		parsigs = append(parsigs, parsig)
	}

	// Create expected aggregated signature
	aggSig, err := tbls.Aggregate(psigs)
	require.NoError(t, err)
	expect, err := aggSig.MarshalBinary()
	require.NoError(t, err)

	agg := sigagg.New(threshold)

	// Assert output
	agg.Subscribe(func(_ context.Context, _ core.Duty, _ core.PubKey, aggData core.AggSignedData) error {
		require.Equal(t, expect, aggData.Signature)

		sig := new(bls_sig.Signature)
		err := sig.UnmarshalBinary(aggData.Signature)
		require.NoError(t, err)

		ok, err := tbls.Verify(tss.PublicKey(), msg, sig)
		require.NoError(t, err)
		require.True(t, ok)

		return nil
	})

	// Run aggregation
	err = agg.Aggregate(ctx, core.Duty{Type: core.DutyAttester}, "", parsigs)
	require.NoError(t, err)
}