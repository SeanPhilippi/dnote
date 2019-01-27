package view

import (
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/dnote/cli/cmd/cat"
	"github.com/dnote/cli/cmd/ls"
	"github.com/dnote/cli/utils"
)

var example = `
 * View all books
 dnote view

 * List notes in a book
 dnote view javascript

 * View a note by an id
 dnote view 1
 `

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) > 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

// NewCmd returns a new view command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <book name?> <note index?>",
		Aliases: []string{"v"},
		Short:   "List books, notes or view a content",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		var run core.RunEFunc

		if len(args) == 0 {
			run = ls.NewRun(ctx)
		} else if len(args) == 1 {
			if utils.IsInt(args[0]) {
				run = cat.NewRun(ctx)
			} else {
				run = ls.NewRun(ctx)
			}
		} else if len(args) == 2 {
			// DEPRECATED: passing book name to view command is deprecated
			run = cat.NewRun(ctx)
		} else {
			return errors.New("Incorrect number of arguments")
		}

		return run(cmd, args)
	}
}
