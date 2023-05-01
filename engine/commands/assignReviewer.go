// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func AssignReviewerCmd() *cobra.Command {
	assignReviewerCmd := &cobra.Command{
		Use:           "assign-reviewers",
		Short:         "Assign reviewers to a pull request",
		Long:          "Assigns a defined amount of reviewers to the pull request from the provided list of reviewers.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}

			if _, err := strconv.ParseUint(args[0], 10, 64); err != nil {
				return fmt.Errorf("invalid argument: %s, number of reviewers must be a number", args[0])
			}

			return nil
		},
		RunE: AssignReviewer,
	}

	return assignReviewerCmd
}

func AssignReviewer(cmd *cobra.Command, args []string) error {
	totalRequiredReviewers, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	if totalRequiredReviewers == 0 {
		return fmt.Errorf("invalid argument: %d, number of reviewers must be greater than 0", totalRequiredReviewers)
	}

	action := fmt.Sprintf(`$assignCodeAuthorReviewers(%d, [], 0)`, totalRequiredReviewers)

	cmd.Print(action)

	return nil
}
