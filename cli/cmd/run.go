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

	"github.com/google/uuid"
	"github.com/reviewpad/api/go/clients"
	github_event_processor "github.com/reviewpad/go-lib/entities/event_processor/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	log "github.com/reviewpad/go-lib/logrus"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v4"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

const (
	CODEHOST_SERVICE_ENDPOINT = "INPUT_CODEHOST_SERVICE"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs reviewpad",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
	runCmd.Flags().BoolVarP(&safeModeRun, "safe-mode-run", "s", false, "Safe mode")
	runCmd.Flags().StringVarP(&url, "url", "u", "", "Code host pull request or issue url")
	runCmd.Flags().StringVarP(&token, "token", "t", "", "Code host token")
	runCmd.Flags().StringVarP(&eventFilePath, "event-payload", "e", "", "File path to github event in JSON format")
	runCmd.Flags().StringVarP(&mixpanelToken, "mixpanel-token", "m", "", "Mixpanel token")
	runCmd.Flags().StringVarP(&logLevel, "log-level", "l", "debug", "Log level")

	if err := runCmd.MarkFlagRequired("url"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("token"); err != nil {
		panic(err)
	}

	if err := runCmd.MarkFlagRequired("event-payload"); err != nil {
		panic(err)
	}

	if endpoint := os.Getenv("INPUT_CODEHOST_SERVICE"); endpoint == "" {
		panic("INPUT_CODEHOST_SERVICE env var is required")
	}

	if endpoint := os.Getenv("INPUT_SEMANTIC_SERVICE"); endpoint == "" {
		panic("INPUT_SEMANTIC_SERVICE env var is required")
	}

	if endpoint := os.Getenv("INPUT_ROBIN_SERVICE"); endpoint == "" {
		panic("INPUT_ROBIN_SERVICE env var is required")
	}
}

func run() error {
	logLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	log := log.NewLogger(logLevel)

	var event *handler.ActionEvent

	eventFileContent, err := os.ReadFile(eventFilePath)
	if err != nil {
		return err
	}

	event, err = parseEvent(string(eventFileContent))
	if err != nil {
		return err
	}

	// Use the provided token
	event.Token = &token

	// FIXME: Abstract to be code host agnostic
	gitHubDetailsRegex := regexp.MustCompile(`github\.com\/(.+)\/(.+)\/(\w+)\/(\d+)`)
	gitHubEntityDetails := gitHubDetailsRegex.FindSubmatch([]byte(url))

	repositoryOwner := string(gitHubEntityDetails[1][:])
	entityKind, err := toTargetEntityKind(string(gitHubEntityDetails[3][:]))
	if err != nil {
		log.Fatalf("Error converting entity kind. Details %+q", err.Error())
	}

	ctx := context.Background()
	gitHubClient := gh.NewGithubClientFromToken(ctx, token)
	collectorClient, err := collector.NewCollector(mixpanelToken, repositoryOwner, string(entityKind), "local-cli", nil)
	if err != nil {
		log.Errorf("error creating new collector: %v", err)
	}

	rawReviewpadFile, err := os.ReadFile(reviewpadFilePath)
	if err != nil {
		return fmt.Errorf("error reading reviewpad file. Details: %v", err.Error())
	}

	bufRviewpadFile := bytes.NewBuffer(rawReviewpadFile)
	reviewpadFile, err := reviewpad.Load(ctx, log, gitHubClient, bufRviewpadFile)
	if err != nil {
		return fmt.Errorf("error running reviewpad team edition. Details %v", err.Error())
	}

	targetEntities, eventDetails, err := github_event_processor.ProcessEvent(log, event, token)
	if err != nil {
		return fmt.Errorf("error processing event. Details %v", err.Error())
	}

	requestID := uuid.New().String()
	md := metadata.Pairs(codehost.RequestIDKey, requestID)
	ctxReq := metadata.NewOutgoingContext(ctx, md)

	codehostClient, codehostConnection, err := clients.NewCodeHostClient(os.Getenv(CODEHOST_SERVICE_ENDPOINT))
	if err != nil {
		return fmt.Errorf("error creating codehost client. Details %v", err.Error())
	}
	defer codehostConnection.Close()

	hostInfo, err := codehost.GetHostInfo(url)
	if err != nil {
		return fmt.Errorf("error getting host info. Details %v", err.Error())
	}

	codeHostClient := &codehost.CodeHostClient{
		Token:          token,
		HostInfo:       hostInfo,
		CodehostClient: codehostClient,
	}

	for _, targetEntity := range targetEntities {
		log.Infof("Processing entity %s/%s#%d", targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		_, _, _, err = reviewpad.Run(ctxReq, log, gitHubClient, codeHostClient, collectorClient, targetEntity, eventDetails, reviewpadFile, nil, dryRun, safeModeRun, nil)
		if err != nil {
			return fmt.Errorf("error running reviewpad team edition. Details %v", err.Error())
		}
	}

	return nil
}

func parseEvent(rawEvent string) (*handler.ActionEvent, error) {
	event := &handler.ActionEvent{}

	err := json.Unmarshal([]byte(rawEvent), event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func toTargetEntityKind(entityType string) (event_processor.TargetEntityKind, error) {
	switch entityType {
	case "issues":
		return event_processor.Issue, nil
	case "pull":
		return event_processor.PullRequest, nil
	default:
		return "", fmt.Errorf("unknown entity type %s", entityType)
	}
}
