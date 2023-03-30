// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RobinSummarizeExtendedCmd() *cobra.Command {
	robinSummarizeExtendedCmd := &cobra.Command{
		Use:           "summarize-extended",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("accepts 0 args, received %d", len(args))
			}
			return nil
		},
		RunE: RobinSummarizeExtended,
	}

	return robinSummarizeExtendedCmd
}

func RobinSummarizeExtended(cmd *cobra.Command, _ []string) error {
	cmd.Print(`$robinSummarize("extended", "openai-gpt-4")`)

	return nil
}
