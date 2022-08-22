// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package target

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/host-event-handler/handler"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
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

func (t CommonTarget) AddAssignees(assignees []string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.AddAssignees(ctx, owner, repo, number, assignees)

	return err
}

func (t CommonTarget) AddLabels(labels []string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.AddLabels(ctx, owner, repo, number, labels)

	return err
}

func (t CommonTarget) GetAvailableAssignees() ([]User, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	assignees := make([]User, 0)

	users, err := t.githubClient.GetIssuesAvailableAssignees(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		assignees = append(assignees, User{
			Login: *user.Login,
		})
	}

	return assignees, nil
}

func (t CommonTarget) Comment(comment string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.CreateComment(ctx, owner, repo, number, &github.IssueComment{Body: github.String(comment)})

	return err
}

func (t CommonTarget) GetComments() ([]Comment, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	cs, err := t.githubClient.GetPullRequestComments(ctx, owner, repo, number, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, len(cs))

	for _, comment := range cs {
		comments = append(comments, Comment{
			Body: *comment.Body,
		})
	}

	return comments, nil
}

func (t CommonTarget) RemoveLabel(labelName string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, err := t.githubClient.RemoveLabelForIssue(ctx, owner, repo, number, labelName)

	return err
}
