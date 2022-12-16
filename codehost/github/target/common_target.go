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

type CommonTarget struct {
	ctx          context.Context
	targetEntity *handler.TargetEntity
	githubClient *gh.GithubClient
}

func NewCommonTarget(ctx context.Context, targetEntity *handler.TargetEntity, githubClient *gh.GithubClient) *CommonTarget {
	return &CommonTarget{
		ctx,
		targetEntity,
		githubClient,
	}
}

func (t *CommonTarget) AddAssignees(assignees []string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.AddAssignees(ctx, owner, repo, number, assignees)

	return err
}

func (t *CommonTarget) AddLabels(labels []string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.AddLabels(ctx, owner, repo, number, labels)

	return err
}

func (t *CommonTarget) GetAvailableAssignees() ([]*codehost.User, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	assignees := make([]*codehost.User, 0)

	users, err := t.githubClient.GetIssuesAvailableAssignees(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		assignees = append(assignees, &codehost.User{
			Login: *user.Login,
		})
	}

	return assignees, nil
}

func (t *CommonTarget) Comment(comment string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.CreateComment(ctx, owner, repo, number, &github.IssueComment{Body: github.String(comment)})

	return err
}

func (t *CommonTarget) GetComments() ([]*codehost.Comment, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	cs, err := t.githubClient.GetComments(ctx, owner, repo, number, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	comments := make([]*codehost.Comment, len(cs))

	for i, comment := range cs {
		comments[i] = &codehost.Comment{
			Body: *comment.Body,
		}
	}

	return comments, nil
}

func (t *CommonTarget) GetProjectFieldsByProjectNumber(projectNumber uint64) ([]*codehost.ProjectField, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	totalRetries := 2

	ghFields, err := t.githubClient.GetProjectFieldsByProjectNumber(ctx, owner, repo, projectNumber, totalRetries)
	if err != nil {
		return nil, err
	}

	fields := make([]*codehost.ProjectField, len(ghFields))
	for i, field := range ghFields {
		fields[i] = &codehost.ProjectField{
			ID:      field.SingleSelectFieldDetails.ID,
			Name:    field.SingleSelectFieldDetails.Name,
			Options: field.SingleSelectFieldDetails.Options,
		}
	}

	return fields, nil
}

func (t *CommonTarget) GetTargetEntity() *handler.TargetEntity {
	return t.targetEntity
}

func (t *CommonTarget) RemoveLabel(labelName string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, err := t.githubClient.RemoveLabelForIssue(ctx, owner, repo, number, labelName)

	// When the label does not exist, the API returns a 404 error.
	// In this case, we ignore the error.
	if err != nil && err.(*github.ErrorResponse).Response.StatusCode == 404 {
		return nil
	}

	return err
}

func (t *CommonTarget) SetProjectFieldSingleSelect(projectItems []gh.GQLProjectV2Item, projectTitle, fieldName, fieldValue string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

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

	fieldDetails := gh.SingleSelectFieldDetails{}
	fieldOptionID := ""

	for _, field := range fields {
		if strings.EqualFold(field.SingleSelectFieldDetails.Name, fieldName) {
			fieldDetails = field.SingleSelectFieldDetails
			break
		}
	}

	if fieldDetails.ID == "" {
		return gh.ErrProjectHasNoSuchField
	}

	for _, option := range fieldDetails.Options {
		if option.Name == fieldValue {
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
