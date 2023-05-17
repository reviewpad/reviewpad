// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"os"

	log "github.com/reviewpad/go-lib/logrus"
	"github.com/reviewpad/reviewpad/v4"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if input reviewpad file is valid",
	RunE: func(cmd *cobra.Command, args []string) error {
		logLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return err
		}

		log := log.NewLogger(logLevel)
		ctx := context.Background()
		githubClient := gh.NewGithubClientFromToken(ctx, token)

		data, err := os.ReadFile(reviewpadFile)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(data)
		reviewpadFile, err := reviewpad.Load(ctx, log, githubClient, buf)
		if err != nil {
			return err
		}

		for _, group := range reviewpadFile.Groups {
			if group.Spec != "" {
				if _, err = aladino.Parse(group.Spec); err != nil {
					return err
				}
			}

			if group.Where != "" {
				if _, err = aladino.Parse(group.Where); err != nil {
					return err
				}
			}
		}

		for _, rules := range reviewpadFile.Rules {
			if rules.Spec != "" {
				if _, err = aladino.Parse(rules.Spec); err != nil {
					return err
				}
			}
		}

		for _, workflow := range reviewpadFile.Workflows {
			for _, rule := range workflow.Rules {
				for _, action := range rule.ExtraActions {
					if _, err = aladino.Parse(action); err != nil {
						return err
					}
				}
			}

			for _, action := range workflow.Actions {
				if _, err = aladino.Parse(action); err != nil {
					return err
				}
			}
		}

		for _, pipeline := range reviewpadFile.Pipelines {
			if pipeline.Trigger != "" {
				if _, err = aladino.Parse(pipeline.Trigger); err != nil {
					return err
				}
			}

			for _, stage := range pipeline.Stages {
				if stage.Until != "" {
					if _, err = aladino.Parse(stage.Until); err != nil {
						return err
					}
				}

				for _, action := range stage.Actions {
					if _, err = aladino.Parse(action); err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
}
