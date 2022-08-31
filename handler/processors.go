// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	reviewpad_gh "github.com/reviewpad/reviewpad/v3/codehost/github"
)

const (
	PullRequest TargetEntityKind = "pull_request"
	Issue       TargetEntityKind = "issue"
)

type TargetEntityKind string

func (entityType TargetEntityKind) String() string {
	switch entityType {
	case Issue:
		return "issues"
	}
	return "pull"
}

type TargetEntity struct {
	Kind   TargetEntityKind
	Number int
	Owner  string
	Repo   string
}

func ParseEvent(rawEvent string) (*ActionEvent, error) {
	event := &ActionEvent{}

	Log("parsing event %v", rawEvent)

	err := json.Unmarshal([]byte(rawEvent), &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func processCronEvent(token string, e *ActionEvent) ([]*TargetEntity, error) {
	Log("processing 'schedule' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()

	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	repoParts := strings.SplitN(*e.Repository, "/", 2)

	owner := repoParts[0]
	repo := repoParts[1]

	issues, _, err := ghClient.ListIssuesByRepo(ctx, owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("get pull requests: %w", err)
	}

	Log("fetched %d issues", len(issues))

	events := make([]*TargetEntity, 0)
	for _, issue := range issues {
		kind := Issue
		if issue.IsPullRequest() {
			kind = PullRequest
		}
		events = append(events, &TargetEntity{
			Kind:   kind,
			Number: *issue.Number,
			Owner:  owner,
			Repo:   repo,
		})
	}

	Log("found events %v", events)

	return events, nil
}

func processIssuesEvent(e *github.IssuesEvent) []*TargetEntity {
	Log("processing 'issues' event")
	Log("found issue %v", *e.Issue.Number)

	return []*TargetEntity{
		{
			Kind:   Issue,
			Number: *e.Issue.Number,
			Owner:  *e.Repo.Owner.Login,
			Repo:   *e.Repo.Name,
		},
	}
}

func processIssueCommentEvent(e *github.IssueCommentEvent) []*TargetEntity {
	Log("processing 'issue_comment' event")
	Log("found issue %v", *e.Issue.Number)

	kind := Issue
	if e.Issue.IsPullRequest() {
		kind = PullRequest
	}

	return []*TargetEntity{
		{
			Kind:   kind,
			Number: *e.Issue.Number,
			Owner:  *e.Repo.Owner.Login,
			Repo:   *e.Repo.Name,
		},
	}
}

func processPullRequestEvent(e *github.PullRequestEvent) []*TargetEntity {
	Log("processing 'pull_request' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
		{
			Kind:   PullRequest,
			Number: *e.PullRequest.Number,
			Owner:  *e.Repo.Owner.Login,
			Repo:   *e.Repo.Name,
		},
	}
}

func processPullRequestReviewEvent(e *github.PullRequestReviewEvent) []*TargetEntity {
	Log("processing 'pull_request_review' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
		{
			Kind:   PullRequest,
			Number: *e.PullRequest.Number,
			Owner:  *e.Repo.Owner.Login,
			Repo:   *e.Repo.Name,
		},
	}
}

func processPullRequestReviewCommentEvent(e *github.PullRequestReviewCommentEvent) []*TargetEntity {
	Log("processing 'pull_request_review_comment' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
		{
			Kind:   PullRequest,
			Number: *e.PullRequest.Number,
			Owner:  *e.Repo.Owner.Login,
			Repo:   *e.Repo.Name,
		},
	}
}

func processPullRequestTargetEvent(e *github.PullRequestTargetEvent) []*TargetEntity {
	Log("processing 'pull_request_target' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
		{
			Kind:   PullRequest,
			Number: *e.PullRequest.Number,
			Owner:  *e.Repo.Owner.Login,
			Repo:   *e.Repo.Name,
		},
	}
}

func processStatusEvent(token string, e *github.StatusEvent) ([]*TargetEntity, error) {
	Log("processing 'status' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()

	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	prs, err := ghClient.GetPullRequests(ctx, *e.Repo.Owner.Login, *e.Repo.Name)
	if err != nil {
		return nil, fmt.Errorf("get pull requests: %w", err)
	}

	Log("fetched %v prs", len(prs))

	for _, pr := range prs {
		if *pr.Head.SHA == *e.SHA {
			Log("found pr %v", *pr.Number)
			return []*TargetEntity{
				{
					Kind:   PullRequest,
					Number: *pr.Number,
					Owner:  *pr.Base.Repo.Owner.Login,
					Repo:   *pr.Base.Repo.Name,
				},
			}, nil
		}
	}

	Log("no pr found with the head sha %v", *e.SHA)

	return []*TargetEntity{}, nil
}

func processWorkflowRunEvent(token string, e *github.WorkflowRunEvent) ([]*TargetEntity, error) {
	Log("processing 'workflow_run' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()
	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	prs, err := ghClient.GetPullRequests(ctx, *e.Repo.Owner.Login, *e.Repo.Name)
	if err != nil {
		return nil, fmt.Errorf("get pull requests: %w", err)
	}

	Log("fetched %v prs", len(prs))

	for _, pr := range prs {
		if *pr.Head.SHA == *e.WorkflowRun.HeadSHA {
			Log("found pr %v", *pr.Number)
			return []*TargetEntity{
				{
					Kind:   PullRequest,
					Number: *pr.Number,
					Owner:  *pr.Base.Repo.Owner.Login,
					Repo:   *pr.Base.Repo.Name,
				},
			}, nil
		}
	}

	Log("no pr found with the head sha %v", *e.WorkflowRun.HeadSHA)

	return []*TargetEntity{}, nil
}

// reviewpad-an: critical
// output: the list of pull requests/issues that are affected by the event.
func ProcessEvent(event *ActionEvent) ([]*TargetEntity, error) {
	// These events do not have an equivalent in the GitHub webhooks, thus
	// parsing them with github.ParseWebhook would return an error.
	// These are the webhook events: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads
	// And these are the "workflow events": https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows
	switch *event.EventName {
	case "schedule":
		return processCronEvent(*event.Token, event)
	}

	eventPayload, err := github.ParseWebHook(*event.EventName, *event.EventPayload)
	if err != nil {
		return nil, fmt.Errorf("parse github webhook: %w", err)
	}

	switch payload := eventPayload.(type) {
	// Handle github events triggered by actions
	// For more information, visit: https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows
	case *github.IssuesEvent:
		return processIssuesEvent(payload), nil
	case *github.IssueCommentEvent:
		return processIssueCommentEvent(payload), nil
	case *github.PullRequestEvent:
		return processPullRequestEvent(payload), nil
	case *github.PullRequestReviewEvent:
		return processPullRequestReviewEvent(payload), nil
	case *github.PullRequestReviewCommentEvent:
		return processPullRequestReviewCommentEvent(payload), nil
	case *github.PullRequestTargetEvent:
		return processPullRequestTargetEvent(payload), nil
	case *github.StatusEvent:
		return processStatusEvent(*event.Token, payload)
	case *github.WorkflowRunEvent:
		return processWorkflowRunEvent(*event.Token, payload)
	}

	return nil, fmt.Errorf("unknown event payload type: %T", eventPayload)
}
