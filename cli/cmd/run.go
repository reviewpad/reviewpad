// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
	runCmd.Flags().BoolVarP(&safeModeRun, "safe-mode-run", "s", false, "Safe mode")
	runCmd.Flags().StringVarP(&pullRequestUrl, "pull-request", "p", "", "GitHub pull request url")
	runCmd.Flags().StringVarP(&gitHubToken, "github-token", "t", "", "GitHub personal access token")
	runCmd.Flags().StringVarP(&eventFilePath, "event-payload", "e", "", "File path to github action event in JSON format")
	runCmd.Flags().StringVarP(&mixpanelToken, "mixpanel-token", "m", "", "Mixpanel token")

	runCmd.MarkFlagRequired("pull-request")
	runCmd.MarkFlagRequired("github-token")
}

type Event struct {
	Payload *json.RawMessage `json:"event,omitempty"`
	Name    *string          `json:"event_name,omitempty"`
}

func parseEvent(rawEvent string) (interface{}, error) {
	ev := &Event{}

	err := json.Unmarshal([]byte(rawEvent), ev)
	if err != nil {
		return nil, err
	}

	return github.ParseWebHook(*ev.Name, *ev.Payload)
}

func run() error {
	var ev interface{}

	if eventFilePath == "" {
		log.Print("[WARN] No event payload provided. Assuming empty event.")
	} else {
		content, err := ioutil.ReadFile(eventFilePath)
		if err != nil {
			return err
		}

		rawEvent := string(content)
		ev, err = parseEvent(rawEvent)
		if err != nil {
			return err
		}
	}

	pullRequestDetailsRegex := regexp.MustCompile(`github\.com\/(.+)\/(.+)\/pull\/(\d+)`)
	pullRequestDetails := pullRequestDetailsRegex.FindSubmatch([]byte(pullRequestUrl))

	repositoryOwner := string(pullRequestDetails[1][:])
	repositoryName := string(pullRequestDetails[2][:])
	pullRequestNumber, err := strconv.Atoi(string(pullRequestDetails[3][:]))
	if err != nil {
		log.Fatalf("Error converting pull request number. Details %+q", err.Error())
	}

	ctx := context.Background()
	githubClient := gh.NewGithubClientFromToken(ctx, gitHubToken)
	collectorClient := collector.NewCollector(mixpanelToken, repositoryOwner)

	ghPullRequest, _, err := githubClient.GetPullRequest(ctx, repositoryOwner, repositoryName, pullRequestNumber)
	if err != nil {
		log.Fatal(err)
	}

	headRepoOwner := *ghPullRequest.Head.Repo.Owner.Login
	headRepoName := *ghPullRequest.Head.Repo.Name
	headRef := *ghPullRequest.Head.Ref

	if *ghPullRequest.Merged {
		if ghPullRequest.Head == nil {
			return fmt.Errorf("team-edition: pull request is merged and head branched was automatically deleted")
		}

		if ghPullRequest.Head.Repo == nil {
			return fmt.Errorf("team-edition: pull request is merged and head repo branched was automatically deleted")
		}

		if ghPullRequest.Head.Repo.DeleteBranchOnMerge == nil || *ghPullRequest.Head.Repo.DeleteBranchOnMerge {
			return fmt.Errorf("team-edition: pull request is merged and head branched was automatically deleted")
		}

		_, _, err := githubClient.GetRepositoryBranch(ctx, headRepoOwner, headRepoName, headRef, true)
		if err != nil {
			return fmt.Errorf("team-edition: pull request is merged and head branched cannot be retrieved")
		}
	}

	data, err := os.ReadFile(reviewpadFile)
	if err != nil {
		return fmt.Errorf("error reading reviewpad file. Details: %v", err.Error())
	}

	buf := bytes.NewBuffer(data)
	file, err := reviewpad.Load(buf)
	if err != nil {
		return fmt.Errorf("error running reviewpad team edition. Details %v", err.Error())
	}

	targetEntity := &handler.TargetEntity{
		Owner:  repositoryOwner,
		Repo:   repositoryName,
		Number: pullRequestNumber,
		Kind:   handler.PullRequest,
	}

	_, err = reviewpad.Run(ctx, githubClient, collectorClient, targetEntity, ev, file, dryRun, safeModeRun)
	if err != nil {
		return fmt.Errorf("error running reviewpad team edition. Details %v", err.Error())
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
