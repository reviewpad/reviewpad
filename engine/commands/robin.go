// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

func RobinCmd() *cobra.Command {
	robinCmd := &cobra.Command{
		Use: "robin",
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

	robinCmd.AddCommand(RobinSummarizeCmd())
	robinCmd.AddCommand(RobinSummarizeExtendedCmd())
	robinCmd.AddCommand(RobinPromptCmd())

	return robinCmd
}
