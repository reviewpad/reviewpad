// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SummarizeCmd() *cobra.Command {
	summarizeCmd := &cobra.Command{
		Use:           "summarize",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("accepts 0 args, received %d", len(args))
			}
			return nil
		},
		RunE: Summarize,
	}

	return summarizeCmd
}

func Summarize(cmd *cobra.Command, _ []string) error {
	cmd.Print(`$robinSummarize("default")`)

	return nil
}
