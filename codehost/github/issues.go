// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v45/github"
)

func (c *GithubClient) CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return c.clientREST.Issues.CreateComment(ctx, owner, repo, number, comment)
}

func (c *GithubClient) EditComment(ctx context.Context, owner string, repo string, commentId int64, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return c.clientREST.Issues.EditComment(ctx, owner, repo, commentId, comment)
}

func (c *GithubClient) DeleteComment(ctx context.Context, owner string, repo string, commentId int64) (*github.Response, error) {
	return c.clientREST.Issues.DeleteComment(ctx, owner, repo, commentId)
}

func (c *GithubClient) CreateLabel(ctx context.Context, owner string, repo string, label *github.Label) (*github.Label, *github.Response, error) {
	return c.clientREST.Issues.CreateLabel(ctx, owner, repo, label)
}

func (c *GithubClient) GetLabel(ctx context.Context, owner string, repo string, name string) (*github.Label, *github.Response, error) {
	return c.clientREST.Issues.GetLabel(ctx, owner, repo, name)
}

func (c *GithubClient) AddLabels(ctx context.Context, owner string, repo string, number int, labels []string) ([]*github.Label, *github.Response, error) {
	return c.clientREST.Issues.AddLabelsToIssue(ctx, owner, repo, number, labels)
}

func (c *GithubClient) RemoveLabelForIssue(ctx context.Context, owner string, repo string, number int, label string) (*github.Response, error) {
	return c.clientREST.Issues.RemoveLabelForIssue(ctx, owner, repo, number, label)
}

func (c *GithubClient) AddAssignees(ctx context.Context, owner string, repo string, number int, assignees []string) (*github.Issue, *github.Response, error) {
	return c.clientREST.Issues.AddAssignees(ctx, owner, repo, number, assignees)
}

func (c *GithubClient) ListIssuesByRepo(ctx context.Context, owner string, repo string, opts *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	return c.clientREST.Issues.ListByRepo(ctx, owner, repo, opts)
}
