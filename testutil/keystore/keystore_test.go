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

package keystore_test

import (
	"os"
	"testing"

	"github.com/dB2510/kryptology/pkg/signatures/bls/bls_sig"
	"github.com/stretchr/testify/require"

	"github.com/obolnetwork/charon/testutil/keystore"
)

func TestStoreLoad(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	var secrets []*bls_sig.SecretKey
	for i := 0; i < 2; i++ {
		_, secret, err := bls_sig.NewSigEth2().Keygen()
		require.NoError(t, err)

		secrets = append(secrets, secret)
	}

	err = keystore.StoreSimnetKeys(secrets, dir)
	require.NoError(t, err)

	actual, err := keystore.LoadSimnetKeys(dir)
	require.NoError(t, err)

	require.Equal(t, secrets, actual)
}