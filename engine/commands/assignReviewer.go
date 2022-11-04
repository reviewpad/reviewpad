// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

			if !regexp.MustCompile(`^([a-zA-Z0-9\-]+[a-zA-Z0-9])(,[a-zA-Z0-9\-]+[a-zA-Z0-9])*$`).MatchString(args[0]) {
				return errors.New("reviewers must be a list of comma separated valid github usernames")
			}

			return nil
		},
		RunE: AssignReviewer,
	}

	flags := assignReviewerCmd.Flags()

	flags.StringP("review-policy", "p", "reviewpad", "The policy followed for reviewer assignment. The valid values can only be: random, round-robin, reviewpad. By default, the policy is reviewpad.")

	flags.Int8P("total-reviewers", "t", 1, "Total number of reviewers, default value is 1")

	return assignReviewerCmd
}

func AssignReviewer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	reviewersList := strings.Split(strings.ReplaceAll(args[0], " ", ""), ",")

	availableReviewers := `"` + strings.Join(reviewersList, `","`) + `"`

	totalRequiredReviewers, err := flags.GetInt8("total-reviewers")
	if err != nil {
		return err
	}

	policy, err := flags.GetString("review-policy")
	if err != nil {
		return err
	}

	if policy != "reviewpad" && policy != "round-robin" && policy != "random" {
		return errors.New("invalid review policy, valid review policies are (reviewpad, round-robin, random)")
	}

	action := fmt.Sprintf(`$assignReviewer([%s], %d, %q)`, availableReviewers, totalRequiredReviewers, policy)

	cmd.Print(action)

	return nil
}
