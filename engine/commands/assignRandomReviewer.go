// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func AssignRandomReviewerCmd() *cobra.Command {
	assignRandomReviewerCmd := &cobra.Command{
		Use:           "assign-random-reviewer",
		Short:         "Assign a random reviewer to a pull request",
		Long:          "Assign a random user from the GitHub organization as a pull request reviewer",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("accepts no args, received %d", len(args))
			}

			return nil
		},
		RunE: AssignRandomReviewer,
	}

	return assignRandomReviewerCmd
}

func AssignRandomReviewer(cmd *cobra.Command, args []string) error {
	action := "$assignRandomReviewer()"

	cmd.Print(action)

	return nil
}
