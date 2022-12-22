// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"log"
	"os"
	"regexp"

	"github.com/reviewpad/reviewpad/v3"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if input reviewpad file is valid",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		githubClient := gh.NewGithubClientFromToken(ctx, gitHubToken)

		githubDetailsRegex := regexp.MustCompile(`github\.com\/(.+)\/(.+)\/(\w+)\/(\d+)`)
		githubEntityDetails := githubDetailsRegex.FindSubmatch([]byte(githubUrl))

		repositoryOwner := string(githubEntityDetails[1][:])
		entityKind, err := toTargetEntityKind(string(githubEntityDetails[3][:]))
		if err != nil {
			log.Fatalf("Error converting entity kind. Details %+q", err.Error())
		}

		collectorClient, err := collector.NewCollector(mixpanelToken, repositoryOwner, string(entityKind), "local-cli", nil)
		if err != nil {
			log.Printf("error creating new collector: %v", err)
		}

		data, err := os.ReadFile(reviewpadFile)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(data)
		reviewpadFile, err := reviewpad.Load(ctx, collectorClient, githubClient, buf)
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
