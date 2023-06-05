// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package target

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/go-github/v52/github"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/shurcooL/githubv4"
)

type IssueTarget struct {
	*CommonTarget

	ctx          context.Context
	targetEntity *entities.TargetEntity
	githubClient *gh.GithubClient
	issue        *github.Issue
}

// ensure IssueTarget conforms to Target interface
var _ codehost.Target = (*IssueTarget)(nil)

func NewIssueTarget(ctx context.Context, targetEntity *entities.TargetEntity, githubClient *gh.GithubClient, issue *github.Issue) *IssueTarget {
	return &IssueTarget{
		NewCommonTarget(ctx, targetEntity, githubClient),
		ctx,
		targetEntity,
		githubClient,
		issue,
	}
}

func (t *IssueTarget) GetNodeID() string {
	return t.issue.GetNodeID()
}

func (t *IssueTarget) Close(comment string, stateReason string) error {
	ctx := t.ctx
	issue := t.issue

	if issue.GetState() == "closed" {
		return nil
	}

	var closeIssueMutation struct {
		CloseIssue struct {
			ClientMutationID string
		} `graphql:"closeIssue(input: $input)"`
	}

	var GQLStateReason githubv4.IssueClosedStateReason
	if stateReason == "completed" {
		GQLStateReason = githubv4.IssueClosedStateReasonCompleted
	} else {
		GQLStateReason = githubv4.IssueClosedStateReasonNotPlanned
	}

	input := githubv4.CloseIssueInput{
		IssueID:     githubv4.ID(issue.GetNodeID()),
		StateReason: &GQLStateReason,
	}

	if err := t.githubClient.GetClientGraphQL().Mutate(ctx, &closeIssueMutation, input, nil); err != nil {
		return err
	}

	if comment != "" {
		if errComment := t.Comment(comment); errComment != nil {
			return errComment
		}
	}

	return nil
}

func (t *IssueTarget) GetLabels() []*pbc.Label {
	issue := t.issue
	labels := make([]*pbc.Label, len(issue.Labels))

	for i, label := range issue.Labels {
		labels[i] = &pbc.Label{
			Id:   label.GetNodeID(),
			Name: label.GetName(),
		}
	}

	return labels
}

func (t *IssueTarget) GetAuthor() (*codehost.User, error) {
	issue := t.issue

	return &codehost.User{
		Login: *issue.User.Login,
	}, nil
}

func (t *IssueTarget) GetProjectByName(name string) (*codehost.Project, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	project, err := t.githubClient.GetProjectV2ByName(ctx, owner, repo, name)
	if err != nil {
		return nil, err
	}

	return &codehost.Project{
		ID:     project.ID,
		Number: project.Number,
	}, nil
}

func (t *IssueTarget) GetAssignees() []*pbc.User {
	issue := t.issue
	assignees := make([]*pbc.User, len(issue.Assignees))

	for i, assignee := range issue.Assignees {
		assignees[i] = &pbc.User{
			Login: *assignee.Login,
		}
	}

	return assignees
}

func (t *IssueTarget) IsLinkedToProject(title string) (bool, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	totalRetries := 2

	projects, err := t.githubClient.GetLinkedProjectsForIssue(ctx, owner, repo, number, totalRetries)
	if err != nil {
		return false, err
	}

	for _, project := range projects {
		if project.Project.Title == title {
			return true, nil
		}
	}

	return false, nil
}

func (t *IssueTarget) GetCommentCount() int64 {
	return int64(t.issue.GetComments())
}

func (t *IssueTarget) GetCreatedAt() string {
	return t.issue.GetCreatedAt().String()
}

func (t *IssueTarget) GetLinkedProjects() ([]gh.GQLProjectV2Item, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	totalRetries := 2

	return t.githubClient.GetLinkedProjectsForIssue(ctx, owner, repo, number, totalRetries)
}

func (t *IssueTarget) GetUpdatedAt() string {
	return t.issue.GetUpdatedAt().String()
}

func (t *IssueTarget) GetDescription() string {
	return t.issue.GetBody()
}

func (t *IssueTarget) GetState() pbc.PullRequestStatus {
	return pbc.PullRequestStatus(pbc.PullRequestStatus_value[strings.ToUpper(t.issue.GetState())])
}

func (t *IssueTarget) GetTitle() string {
	return t.issue.GetTitle()
}

func (t *IssueTarget) JSON() (string, error) {
	j, err := json.Marshal(t.issue)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func (target *IssueTarget) GetProjectV2ItemID(projectID string) (string, error) {
	ctx := target.ctx
	targetEntity := target.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	return target.githubClient.GetIssueProjectV2ItemID(ctx, owner, repo, projectID, targetEntity.Number)
}

func (t *IssueTarget) IsInProject(projectTitle string) (bool, error) {
	projectItems, err := t.GetLinkedProjects()
	if err != nil {
		return false, err
	}

	for _, projectItem := range projectItems {
		if projectItem.Project.Title == projectTitle {
			return true, nil
		}
	}

	return false, nil
}

func (t *IssueTarget) AddToProject(projectID string) (string, error) {
	var addProjectV2ItemByIdMutation struct {
		AddProjectV2ItemById struct {
			Item struct {
				Id string
			}
		} `graphql:"addProjectV2ItemById(input: $input)"`
	}

	input := gh.AddProjectV2ItemByIdInput{
		ProjectID: projectID,
		ContentID: t.issue.GetNodeID(),
	}

	err := t.githubClient.GetClientGraphQL().Mutate(t.ctx, &addProjectV2ItemByIdMutation, input, nil)
	if err != nil {
		return "", err
	}

	return addProjectV2ItemByIdMutation.AddProjectV2ItemById.Item.Id, nil
}
