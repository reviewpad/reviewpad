// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func DryRunCmd() *cobra.Command {
	dryRunCmd := &cobra.Command{
		Use:           "dry-run",
		Short:         "Runs the reviewpad configuration in dry-run mode",
		Long:          "Runs the reviewpad configuration in the root repository in dry-run mode and adds the result as a comment.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("accepts no args, received %d", len(args))
			}

			return nil
		},
		RunE: DryRun,
	}

	return dryRunCmd
}

func DryRun(cmd *cobra.Command, args []string) error {
	action := "$dryRun()"

	cmd.Print(action)

	return nil
}
