// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
)

func NewCommands(out io.Writer, args []string) *cobra.Command {
	root := &cobra.Command{
		Use:    "/reviewpad",
		Hidden: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Lookup("help").Hidden = true
			return errors.New(cmd.UsageString())
		},
		DisableFlagParsing: true,
	}

	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.SetOut(out)

	root.SetArgs(args)

	root.AddCommand(AssignReviewerCmd())
	root.AddCommand(AssignRandomReviewerCmd())

	return root
}
