// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package target

import (
	"context"
	"errors"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/host-event-handler/handler"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
)

var (
	ErrReviewNotSupported = errors.New("review not supported on issues")
	ErrMergeUnsupported   = errors.New("merge not supported")
	ErrCommitNotSupported = errors.New("commit not supported")
	ErrFileNotSupported   = errors.New("file not supported")
	ErrIssuesNotSupported = errors.New("issues not supported")
	ErrRefNotSupported    = errors.New("reference not supported")
	ErrDraftNotSupported  = errors.New("draft not supported")
)

type IssueTarget struct {
	*CommonTarget

	ctx          context.Context
	targetEntity *handler.TargetEntity
	githubClient *gh.GithubClient
	issue        *github.Issue
}

// ensure IssueTarget conforms to Target interface
var _ Target = IssueTarget{}

func NewIssueTarget(ctx context.Context, targetEntity *handler.TargetEntity, githubClient *gh.GithubClient, issue *github.Issue) *IssueTarget {
	return &IssueTarget{
		NewCommonTarget(ctx, targetEntity, githubClient),
		ctx,
		targetEntity,
		githubClient,
		issue,
	}
}

func (t IssueTarget) AddToProject(projectID, fieldID, optionID string) error {
	// ctx := t.env.GetCtx()
	// pr := t.env.GetPullRequest()

	// return t.env.GetGithubClient().AddToProject(ctx, projectID, *pr.NodeID, fieldID, optionID)
	return nil
}

func (t IssueTarget) Close() error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	issue := t.issue
	issue.State = github.String("closed")
	issueRequest := &github.IssueRequest{
		State: issue.State,
	}

	_, _, err := t.githubClient.EditIssue(ctx, owner, repo, number, issueRequest)

	return err
}

func (t IssueTarget) GetLabels() ([]Label, error) {
	issue := t.issue
	labels := make([]Label, len(issue.Labels))

	for _, label := range issue.Labels {
		labels = append(labels, Label{
			ID:   *label.ID,
			Name: *label.Name,
		})
	}

	return labels, nil
}

func (t IssueTarget) GetProjectByName(name string) (*Project, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	project, err := t.githubClient.GetProjectV2ByName(ctx, owner, repo, name)
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:     project.ID,
		Number: project.Number,
	}, nil
}

func (t IssueTarget) GetProjectFieldsByProjectNumber(projectNumber uint64) ([]ProjectField, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	totalRetries := 2

	ghFields, err := t.githubClient.GetProjectFieldsByProjectNumber(ctx, owner, repo, projectNumber, totalRetries)
	if err != nil {
		return nil, err
	}

	fields := make([]ProjectField, len(ghFields))
	for _, field := range ghFields {
		fields = append(fields, ProjectField{
			ID:      field.Details.ID,
			Name:    field.Details.Name,
			Options: field.Details.Options,
		})
	}

	return fields, nil
}

func (t IssueTarget) GetRequestedReviewers() ([]User, error) {
	return nil, ErrReviewNotSupported
}

func (t IssueTarget) GetReviewers() (*Reviewers, error) {
	return nil, ErrReviewNotSupported
}

func (t IssueTarget) GetReviews() ([]Review, error) {
	return nil, ErrReviewNotSupported
}

func (t IssueTarget) GetUser() (*User, error) {
	issue := t.issue

	return &User{
		Login: *issue.User.Login,
	}, nil
}

func (t IssueTarget) Merge(mergeMethod string) error {
	return ErrMergeUnsupported
}

func (t IssueTarget) RequestReviewers(reviewers []string) error {
	return ErrReviewNotSupported
}

func (t IssueTarget) RequestTeamReviewers(reviewers []string) error {
	return ErrReviewNotSupported
}

func (t IssueTarget) GetAssignees() ([]User, error) {
	issue := t.issue
	assignees := make([]User, len(issue.Assignees))

	for _, assignee := range issue.Assignees {
		assignees = append(assignees, User{
			Login: *assignee.Login,
		})
	}

	return assignees, nil
}

func (t IssueTarget) GetBase() (string, error) {
	return "", ErrRefNotSupported
}

func (t IssueTarget) GetCommentCount() (int, error) {
	return t.issue.GetComments(), nil
}

func (t IssueTarget) GetCommitCount() (int, error) {
	return 0, ErrCommitNotSupported
}

func (t IssueTarget) GetCommits() ([]Commit, error) {
	return nil, ErrCommitNotSupported
}

func (t IssueTarget) GetCreatedAt() (string, error) {
	return t.issue.GetCreatedAt().String(), nil
}

func (t IssueTarget) GetDescription() (string, error) {
	return t.issue.GetBody(), nil
}

func (t IssueTarget) GetFileCount() (int, error) {
	return 0, ErrFileNotSupported
}

func (t IssueTarget) GetLinkedIssuesCount() (int, error) {
	return 0, ErrIssuesNotSupported
}

func (t IssueTarget) GetReviewThreads() ([]ReviewThread, error) {
	return nil, ErrReviewNotSupported
}

func (t IssueTarget) GetHead() (string, error) {
	return "", ErrRefNotSupported
}

func (t IssueTarget) IsDraft() (bool, error) {
	return false, nil
}
