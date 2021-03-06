/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package logout

import (
	"database/sql"

	"github.com/dnote/dnote/cli/client"
	"github.com/dnote/dnote/cli/core"
	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ErrNotLoggedIn is an error for logging out when not logged in
var ErrNotLoggedIn = errors.New("not logged in")

var example = `
  dnote logout`

// NewCmd returns a new logout command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Logout from the server",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

// Do performs logout
func Do(ctx infra.DnoteCtx) error {
	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	var key string
	err = core.GetSystem(tx, infra.SystemSessionKey, &key)
	if errors.Cause(err) == sql.ErrNoRows {
		return ErrNotLoggedIn
	} else if err != nil {
		return errors.Wrap(err, "getting session key")
	}

	err = client.Signout(ctx, key)
	if err != nil {
		return errors.Wrap(err, "requesting logout")
	}

	if err := core.DeleteSystem(tx, infra.SystemCipherKey); err != nil {
		return errors.Wrap(err, "deleting enc key")
	}
	if err := core.DeleteSystem(tx, infra.SystemSessionKey); err != nil {
		return errors.Wrap(err, "deleting session key")
	}
	if err := core.DeleteSystem(tx, infra.SystemSessionKeyExpiry); err != nil {
		return errors.Wrap(err, "deleting session key expiry")
	}

	tx.Commit()

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		err := Do(ctx)
		if err == ErrNotLoggedIn {
			log.Error("not logged in\n")
			return nil
		} else if err != nil {
			return errors.Wrap(err, "logging out")
		}

		log.Success("logged out\n")

		return nil
	}
}
