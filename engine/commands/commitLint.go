// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CommitLintCmd() *cobra.Command {
	commitLintCmd := &cobra.Command{
		Use:           "commit-lint",
		Short:         "Check if commit messages meet the conventional commit format.",
		Long:          "Checks if the commits in the pull request follow the conventional commits specification.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("accepts no args, received %d", len(args))
			}

			return nil
		},
		RunE: CommitLint,
	}

	return commitLintCmd
}

func CommitLint(cmd *cobra.Command, args []string) error {
	action := "$commitLint()"

	cmd.Print(action)

	return nil
}
