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

package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/p2p"
)

func newCreateEnrCmd(runFunc func(io.Writer, p2p.Config, string) error) *cobra.Command {
	var (
		config  p2p.Config
		dataDir string
	)

	cmd := &cobra.Command{
		Use:   "enr",
		Short: "Create an Ethereum Node Record (ENR) private key to identify this charon client",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFunc(cmd.OutOrStdout(), config, dataDir)
		},
	}

	bindDataDirFlag(cmd.Flags(), &dataDir)
	bindP2PFlags(cmd.Flags(), &config)

	return cmd
}

// runCreateEnrCmd stores a new p2pkey to disk and prints the ENR for the provided config.
func runCreateEnrCmd(w io.Writer, config p2p.Config, dataDir string) error {
	key, err := p2p.NewSavedPrivKey(dataDir)
	if err != nil {
		return err
	}

	localEnode, db, err := p2p.NewLocalEnode(config, key)
	if err != nil {
		return errors.Wrap(err, "failed to open peer DB")
	}
	defer db.Close()

	_, _ = fmt.Fprintf(w, "Created ENR private key: %s/p2pkey\n", dataDir)
	_, _ = fmt.Fprintln(w, localEnode.Node().String())

	return nil
}
