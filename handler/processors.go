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

	"github.com/google/go-github/v48/github"
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

type EventData struct {
	EventName   string
	EventAction string
	Comment     *github.IssueComment
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

func processUnsupportedEvent(eventPayload interface{}) ([]*TargetEntity, []*EventData, error) {
	return nil, nil, fmt.Errorf("unsupported event payload type: %T", eventPayload)
}

func processCronEvent(token string, e *ActionEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'schedule' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()

	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	repoParts := strings.SplitN(*e.Repository, "/", 2)

	owner := repoParts[0]
	repo := repoParts[1]

	issues, _, err := ghClient.ListIssuesByRepo(ctx, owner, repo, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("get pull requests: %w", err)
	}

	Log("fetched %d issues", len(issues))

	targets := make([]*TargetEntity, 0)
	events := make([]*EventData, 0)
	for _, issue := range issues {
		kind := Issue
		if issue.IsPullRequest() {
			kind = PullRequest
		}
		targets = append(targets, &TargetEntity{
			Kind:   kind,
			Number: *issue.Number,
			Owner:  owner,
			Repo:   repo,
		})

		events = append(events, &EventData{
			EventName: *e.EventName,
		})
	}

	Log("found events %v", targets)

	return targets, events, nil
}

func processIssuesEvent(e *github.IssuesEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'issues' event")
	Log("found issue %v", *e.Issue.Number)

	return []*TargetEntity{
			{
				Kind:   Issue,
				Number: *e.Issue.Number,
				Owner:  *e.Repo.Owner.Login,
				Repo:   *e.Repo.Name,
			},
		}, []*EventData{
			{
				EventName:   "issues",
				EventAction: *e.Action,
			},
		}, nil
}

func processIssueCommentEvent(e *github.IssueCommentEvent) ([]*TargetEntity, []*EventData, error) {
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
		}, []*EventData{
			{
				EventName:   "issue_comment",
				Comment:     e.Comment,
				EventAction: *e.Action,
			},
		}, nil
}

func processPullRequestEvent(e *github.PullRequestEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'pull_request' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
			{
				Kind:   PullRequest,
				Number: *e.PullRequest.Number,
				Owner:  *e.Repo.Owner.Login,
				Repo:   *e.Repo.Name,
			},
		}, []*EventData{
			{
				EventName:   "pull_request",
				EventAction: *e.Action,
			},
		}, nil
}

func processPullRequestReviewEvent(e *github.PullRequestReviewEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'pull_request_review' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
			{
				Kind:   PullRequest,
				Number: *e.PullRequest.Number,
				Owner:  *e.Repo.Owner.Login,
				Repo:   *e.Repo.Name,
			},
		}, []*EventData{
			{
				EventName:   "pull_request_review",
				EventAction: *e.Action,
			},
		}, nil
}

func processPullRequestReviewCommentEvent(e *github.PullRequestReviewCommentEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'pull_request_review_comment' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
			{
				Kind:   PullRequest,
				Number: *e.PullRequest.Number,
				Owner:  *e.Repo.Owner.Login,
				Repo:   *e.Repo.Name,
			},
		}, []*EventData{
			{
				EventName:   "pull_request_review_comment",
				EventAction: *e.Action,
			},
		}, nil
}

func processPullRequestTargetEvent(e *github.PullRequestTargetEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'pull_request_target' event")
	Log("found pr %v", *e.PullRequest.Number)

	return []*TargetEntity{
			{
				Kind:   PullRequest,
				Number: *e.PullRequest.Number,
				Owner:  *e.Repo.Owner.Login,
				Repo:   *e.Repo.Name,
			},
		}, []*EventData{
			{
				EventName:   "pull_request_target",
				EventAction: *e.Action,
			},
		}, nil
}

func processStatusEvent(token string, e *github.StatusEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'status' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()

	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	prs, err := ghClient.GetPullRequests(ctx, *e.Repo.Owner.Login, *e.Repo.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("get pull requests: %w", err)
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
				}, []*EventData{
					{
						EventName: "status",
					},
				}, nil
		}
	}

	Log("no pr found with the head sha %v", *e.SHA)

	return []*TargetEntity{}, []*EventData{}, nil
}

func processWorkflowRunEvent(token string, e *github.WorkflowRunEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'workflow_run' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()
	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	prs, err := ghClient.GetPullRequests(ctx, *e.Repo.Owner.Login, *e.Repo.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("get pull requests: %w", err)
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
				}, []*EventData{
					{
						EventName:   "workflow_run",
						EventAction: *e.Action,
					},
				}, nil
		}
	}

	Log("no pr found with the head sha %v", *e.WorkflowRun.HeadSHA)

	return []*TargetEntity{}, []*EventData{}, nil
}

func processPushEvent(token string, e *github.PushEvent) ([]*TargetEntity, []*EventData, error) {
	Log("processing 'push' event")

	ctx, canc := context.WithTimeout(context.Background(), time.Minute*10)
	defer canc()

	ghClient := reviewpad_gh.NewGithubClientFromToken(ctx, token)

	repoParts := strings.SplitN(*e.GetRepo().FullName, "/", 2)

	owner := repoParts[0]
	repo := repoParts[1]

	prs, err := ghClient.GetPullRequests(ctx, owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("get pull requests: %w", err)
	}

	Log("fetched %d pull requests", len(prs))

	targets := make([]*TargetEntity, 0)
	events := make([]*EventData, 0)
	for _, pr := range prs {
		if pr.Base.GetRef() == e.GetRef() {
			targets = append(targets, &TargetEntity{
				Kind:   PullRequest,
				Number: *pr.Number,
				Owner:  owner,
				Repo:   repo,
			})
			events = append(events, &EventData{
				EventName: "push",
			})
		}
	}

	Log("found events %v", targets)

	return targets, events, nil
}

// reviewpad-an: critical
// output: the list of pull requests/issues that are affected by the event.
func ProcessEvent(event *ActionEvent) ([]*TargetEntity, []*EventData, error) {
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
		return nil, nil, fmt.Errorf("parse github webhook: %w", err)
	}

	switch payload := eventPayload.(type) {
	// Handle github events triggered by actions
	// For more information, visit: https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows
	case *github.BranchProtectionRuleEvent:
		return processUnsupportedEvent(payload)
	case *github.CheckRunEvent:
		return processUnsupportedEvent(payload)
	case *github.CheckSuiteEvent:
		return processUnsupportedEvent(payload)
	case *github.CommitCommentEvent:
		return processUnsupportedEvent(payload)
	case *github.ContentReferenceEvent:
		return processUnsupportedEvent(payload)
	case *github.CreateEvent:
		return processUnsupportedEvent(payload)
	case *github.DeleteEvent:
		return processUnsupportedEvent(payload)
	case *github.DeployKeyEvent:
		return processUnsupportedEvent(payload)
	case *github.DeploymentEvent:
		return processUnsupportedEvent(payload)
	case *github.DeploymentStatusEvent:
		return processUnsupportedEvent(payload)
	case *github.DiscussionEvent:
		return processUnsupportedEvent(payload)
	case *github.ForkEvent:
		return processUnsupportedEvent(payload)
	case *github.GitHubAppAuthorizationEvent:
		return processUnsupportedEvent(payload)
	case *github.GollumEvent:
		return processUnsupportedEvent(payload)
	case *github.InstallationEvent:
		return processUnsupportedEvent(payload)
	case *github.InstallationRepositoriesEvent:
		return processUnsupportedEvent(payload)
	case *github.IssueCommentEvent:
		return processIssueCommentEvent(payload)
	case *github.IssuesEvent:
		return processIssuesEvent(payload)
	case *github.LabelEvent:
		return processUnsupportedEvent(payload)
	case *github.MarketplacePurchaseEvent:
		return processUnsupportedEvent(payload)
	case *github.MemberEvent:
		return processUnsupportedEvent(payload)
	case *github.MembershipEvent:
		return processUnsupportedEvent(payload)
	case *github.MetaEvent:
		return processUnsupportedEvent(payload)
	case *github.MilestoneEvent:
		return processUnsupportedEvent(payload)
	case *github.OrganizationEvent:
		return processUnsupportedEvent(payload)
	case *github.OrgBlockEvent:
		return processUnsupportedEvent(payload)
	case *github.PackageEvent:
		return processUnsupportedEvent(payload)
	case *github.PageBuildEvent:
		return processUnsupportedEvent(payload)
	case *github.PingEvent:
		return processUnsupportedEvent(payload)
	case *github.ProjectEvent:
		return processUnsupportedEvent(payload)
	case *github.ProjectCardEvent:
		return processUnsupportedEvent(payload)
	case *github.ProjectColumnEvent:
		return processUnsupportedEvent(payload)
	case *github.PublicEvent:
		return processUnsupportedEvent(payload)
	case *github.PullRequestEvent:
		return processPullRequestEvent(payload)
	case *github.PullRequestReviewEvent:
		return processPullRequestReviewEvent(payload)
	case *github.PullRequestReviewCommentEvent:
		return processPullRequestReviewCommentEvent(payload)
	case *github.PullRequestReviewThreadEvent:
		return processUnsupportedEvent(payload)
	case *github.PullRequestTargetEvent:
		return processPullRequestTargetEvent(payload)
	case *github.PushEvent:
		return processPushEvent(*event.Token, payload)
	case *github.ReleaseEvent:
		return processUnsupportedEvent(payload)
	case *github.RepositoryEvent:
		return processUnsupportedEvent(payload)
	case *github.RepositoryDispatchEvent:
		return processUnsupportedEvent(payload)
	case *github.RepositoryImportEvent:
		return processUnsupportedEvent(payload)
	case *github.RepositoryVulnerabilityAlertEvent:
		return processUnsupportedEvent(payload)
	case *github.SecretScanningAlertEvent:
		return processUnsupportedEvent(payload)
	case *github.StarEvent:
		return processUnsupportedEvent(payload)
	case *github.StatusEvent:
		return processStatusEvent(*event.Token, payload)
	case *github.TeamEvent:
		return processUnsupportedEvent(payload)
	case *github.TeamAddEvent:
		return processUnsupportedEvent(payload)
	case *github.UserEvent:
		return processUnsupportedEvent(payload)
	case *github.WatchEvent:
		return processUnsupportedEvent(payload)
	case *github.WorkflowDispatchEvent:
		return processUnsupportedEvent(payload)
	case *github.WorkflowJobEvent:
		return processUnsupportedEvent(payload)
	case *github.WorkflowRunEvent:
		return processWorkflowRunEvent(*event.Token, payload)
	}

	return processUnsupportedEvent(eventPayload)
}
