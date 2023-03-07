// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func RobinCmd() *cobra.Command {
	robinCmd := &cobra.Command{
		Use:           "robin",
		Short:         "Sneak peak at pandora box",
		Long:          "Sneak peak at pandora box",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: Robin,
	}

	return robinCmd
}

func Robin(cmd *cobra.Command, args []string) error {
	prompt := strings.Join(args, " ")

	action := fmt.Sprintf(`$robin("%s")`, prompt)

	cmd.Print(action)

	return nil
}
