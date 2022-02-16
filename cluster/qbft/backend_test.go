package qbft

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/corverroos/quorum/log"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestBackend(t *testing.T) {

	const n = 3
	seeds := []int64{2, 8, 1} // results in addresses 0x0.., 0x1.., 0x2..

	tr := NewMemTransport()

	log.Root().SetHandler(log.FuncHandler(func(r *log.Record) error {
		fmt.Printf("%s %v %v %s %v\n", r.Time.Format("15:04:05.000"), r.Lvl, r.Call, r.Msg, r.Ctx)
		return nil
	}))

	var privkeys []*ecdsa.PrivateKey
	var pubkeys []ecdsa.PublicKey
	for i := 0; i < n; i++ {
		key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.New(rand.NewSource(seeds[i])))
		require.NoError(t, err)

		privkeys = append(privkeys, key)
		pubkeys = append(pubkeys, key.PublicKey)
	}

	var bcks []backend
	for i := 0; i < n; i++ {
		b := New(tr, i, pubkeys, privkeys[i])

		require.NoError(t, b.Start())

		bcks = append(bcks, b)
	}

	for {
		time.Sleep(time.Second)
	}
}
