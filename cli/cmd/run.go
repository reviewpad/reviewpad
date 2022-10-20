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
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
	runCmd.Flags().BoolVarP(&safeModeRun, "safe-mode-run", "s", false, "Safe mode")
	runCmd.Flags().StringVarP(&githubUrl, "github-url", "u", "", "GitHub pull request or issue url")
	runCmd.Flags().StringVarP(&gitHubToken, "github-token", "t", "", "GitHub personal access token")
	runCmd.Flags().StringVarP(&eventFilePath, "event-payload", "e", "", "File path to github action event in JSON format")
	runCmd.Flags().StringVarP(&mixpanelToken, "mixpanel-token", "m", "", "Mixpanel token")

	runCmd.MarkFlagRequired("github-url")
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

func getGithubData(accessToken string) *github.User {
	// Get request to a set URL
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	var actionActor *github.User
	err = json.Unmarshal(respbody, &actionActor)
	if err != nil {
		log.Panic("Could not get user")
	}

	return actionActor
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

	githubDetailsRegex := regexp.MustCompile(`github\.com\/(.+)\/(.+)\/(\w+)\/(\d+)`)
	githubEntityDetails := githubDetailsRegex.FindSubmatch([]byte(githubUrl))

	repositoryOwner := string(githubEntityDetails[1][:])
	repositoryName := string(githubEntityDetails[2][:])
	entityKind, err := toTargetEntityKind(string(githubEntityDetails[3][:]))
	if err != nil {
		log.Fatalf("Error converting entity kind. Details %+q", err.Error())
	}

	entityNumber, err := strconv.Atoi(string(githubEntityDetails[4][:]))
	if err != nil {
		log.Fatalf("Error converting entity number. Details %+q", err.Error())
	}

	ctx := context.Background()
	githubActionActor := getGithubData(gitHubToken)
	githubClient := gh.NewGithubClientFromToken(ctx, gitHubToken)
	collectorClient := collector.NewCollector(mixpanelToken, repositoryOwner, string(entityKind), githubUrl, "local-cli")

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
		Number: entityNumber,
		Kind:   entityKind,
	}

	_, err = reviewpad.Run(ctx, githubActionActor, githubClient, collectorClient, targetEntity, ev, file, dryRun, safeModeRun)
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
