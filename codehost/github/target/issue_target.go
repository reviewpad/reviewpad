// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package target

import (
	"context"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/codehost"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
)

type IssueTarget struct {
	*CommonTarget

	ctx          context.Context
	targetEntity *handler.TargetEntity
	githubClient *gh.GithubClient
	issue        *github.Issue
}

// ensure IssueTarget conforms to Target interface
var _ codehost.Target = (*IssueTarget)(nil)

func NewIssueTarget(ctx context.Context, targetEntity *handler.TargetEntity, githubClient *gh.GithubClient, issue *github.Issue) *IssueTarget {
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
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	issue := t.issue
	issue.State = github.String("closed")
	issueRequest := &github.IssueRequest{
		State:       issue.State,
		StateReason: github.String(stateReason),
	}

	_, _, err := t.githubClient.EditIssue(ctx, owner, repo, number, issueRequest)
	if err != nil {
		return err
	}

	if comment != "" {
		if errComment := t.Comment(comment); errComment != nil {
			return errComment
		}
	}

	return err
}

func (t *IssueTarget) GetLabels() []*codehost.Label {
	issue := t.issue
	labels := make([]*codehost.Label, len(issue.Labels))

	for i, label := range issue.Labels {
		labels[i] = &codehost.Label{
			ID:   *label.ID,
			Name: *label.Name,
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

func (t *IssueTarget) GetAssignees() ([]*codehost.User, error) {
	issue := t.issue
	assignees := make([]*codehost.User, len(issue.Assignees))

	for i, assignee := range issue.Assignees {
		assignees[i] = &codehost.User{
			Login: *assignee.Login,
		}
	}

	return assignees, nil
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

func (t *IssueTarget) GetCommentCount() (int, error) {
	return t.issue.GetComments(), nil
}

func (t *IssueTarget) GetCreatedAt() (string, error) {
	return t.issue.GetCreatedAt().String(), nil
}

func (t *IssueTarget) GetUpdatedAt() (string, error) {
	return t.issue.GetUpdatedAt().String(), nil
}

func (t *IssueTarget) GetDescription() (string, error) {
	return t.issue.GetBody(), nil
}

func (t *IssueTarget) GetState() string {
	return t.issue.GetState()
}

func (t *IssueTarget) GetTitle() string {
	return t.issue.GetTitle()
}

func (t *IssueTarget) SetProjectFieldSingleSelect(projectTitle string, fieldName string, fieldValue string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	totalRetries := 2

	projectItems, err := t.githubClient.GetLinkedProjectsForIssue(ctx, owner, repo, number, totalRetries)
	if err != nil {
		return err
	}

	var projectItemID string
	var projectID string
	var projectNumber uint64
	foundProject := false
	totalRequestTries := 2

	for _, projectItem := range projectItems {
		if projectItem.Project.Title == projectTitle {
			projectItemID = projectItem.ID
			projectNumber = projectItem.Project.Number
			projectID = projectItem.Project.ID
			foundProject = true
			break
		}
	}

	if !foundProject {
		return gh.ErrProjectNotFound
	}

	fields, err := t.githubClient.GetProjectFieldsByProjectNumber(ctx, owner, repo, projectNumber, totalRequestTries)
	if err != nil {
		return err
	}

	fieldDetails := gh.FieldDetails{}
	fieldOptionID := ""

	for _, field := range fields {
		if strings.EqualFold(field.Details.Name, fieldName) {
			fieldDetails = field.Details
			break
		}
	}

	if fieldDetails.ID == "" {
		return gh.ErrProjectHasNoSuchField
	}

	for _, option := range fieldDetails.Options {
		if strings.Contains(strings.ToLower(option.Name), fieldValue) {
			fieldOptionID = option.ID
			break
		}
	}

	if fieldOptionID == "" {
		return gh.ErrProjectHasNoSuchFieldValue
	}

	var updateProjectV2ItemFieldValueMutation struct {
		UpdateProjetV2ItemFieldValue struct {
			ClientMutationID string
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	updateInput := gh.UpdateProjectV2ItemFieldValueInput{
		ProjectID: projectID,
		ItemID:    projectItemID,
		Value: gh.FieldValue{
			SingleSelectOptionId: fieldOptionID,
		},
		FieldID: fieldDetails.ID,
	}

	return t.githubClient.GetClientGraphQL().Mutate(ctx, &updateProjectV2ItemFieldValueMutation, updateInput, nil)
}
