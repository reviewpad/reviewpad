// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/reviewpad/reviewpad/v3"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
	runCmd.Flags().BoolVarP(&safeModeRun, "safe-mode-run", "s", false, "Safe mode")
	runCmd.Flags().StringVarP(&gitHubUrl, "github-url", "u", "", "GitHub pull request or issue url")
	runCmd.Flags().StringVarP(&gitHubToken, "github-token", "t", "", "GitHub personal access token")
	runCmd.Flags().StringVarP(&eventFilePath, "event-payload", "e", "", "File path to github event in JSON format")
	runCmd.Flags().StringVarP(&mixpanelToken, "mixpanel-token", "m", "", "Mixpanel token")
	runCmd.Flags().StringVarP(&logLevel, "log-level", "l", "debug", "Log level")

	if err := runCmd.MarkFlagRequired("github-url"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("github-token"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("event-payload"); err != nil {
		panic(err)
	}
}

func parseEvent(rawEvent string) (*handler.ActionEvent, error) {
	ev := &handler.ActionEvent{}

	err := json.Unmarshal([]byte(rawEvent), ev)
	if err != nil {
		return nil, err
	}

	return ev, nil
}

func toTargetEntityKind(entityType string) (handler.TargetEntityKind, error) {
	switch entityType {
	case "issues":
		return handler.Issue, nil
	case "pull":
		return handler.PullRequest, nil
	default:
		return "", fmt.Errorf("unknown entity type %s", entityType)
	}
}

func run() error {
	logLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	log := utils.NewLogger(logLevel)

	var event *handler.ActionEvent

	if eventFilePath == "" {
		log.Warn("[WARN] No event payload provided. Assuming empty event.")
	} else {
		content, readFileErr := os.ReadFile(eventFilePath)
		if readFileErr != nil {
			return readFileErr
		}

		rawEvent := string(content)
		event, err = parseEvent(rawEvent)
		if err != nil {
			return err
		}
	}

	gitHubDetailsRegex := regexp.MustCompile(`github\.com\/(.+)\/(.+)\/(\w+)\/(\d+)`)
	gitHubEntityDetails := gitHubDetailsRegex.FindSubmatch([]byte(gitHubUrl))

	repositoryOwner := string(gitHubEntityDetails[1][:])
	entityKind, err := toTargetEntityKind(string(gitHubEntityDetails[3][:]))
	if err != nil {
		log.Fatalf("Error converting entity kind. Details %+q", err.Error())
	}

	ctx := context.Background()
	gitHubClient := gh.NewGithubClientFromToken(ctx, gitHubToken)
	collectorClient, err := collector.NewCollector(mixpanelToken, repositoryOwner, string(entityKind), "local-cli", nil)
	if err != nil {
		log.Errorf("error creating new collector: %v", err)
	}

	data, err := os.ReadFile(reviewpadFile)
	if err != nil {
		return fmt.Errorf("error reading reviewpad file. Details: %v", err.Error())
	}

	buf := bytes.NewBuffer(data)
	file, err := reviewpad.Load(ctx, log, gitHubClient, buf)
	if err != nil {
		return fmt.Errorf("error running reviewpad team edition. Details %v", err.Error())
	}

	targetEntities, eventDetails, err := handler.ProcessEvent(log, event)
	if err != nil {
		return fmt.Errorf("error processing event. Details %v", err.Error())
	}

	for _, targetEntity := range targetEntities {
		log.Infof("Processing entity %s/%s#%d", targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		_, _, err = reviewpad.Run(ctx, log, gitHubClient, collectorClient, targetEntity, eventDetails, file, dryRun, safeModeRun)
		if err != nil {
			return fmt.Errorf("error running reviewpad team edition. Details %v", err.Error())
		}
	}

	return nil
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs reviewpad",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}
