// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v52/github"
	"github.com/shurcooL/githubv4"
)

type GQLProjectV2Item struct {
	ID      string
	Project ProjectV2
}

type IssueLinkedProjectsQuery struct {
	Repository struct {
		Issue struct {
			ProjectItems *struct {
				Nodes    []GQLProjectV2Item
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"projectItems(first: 10, after: $projectItemsCursor)"`
		} `graphql:"issue(number: $issueNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

func (c *GithubClient) GetIssue(ctx context.Context, owner, repo string, number int) (*github.Issue, *github.Response, error) {
	return c.clientREST.Issues.Get(ctx, owner, repo, number)
}

func (c *GithubClient) EditIssue(ctx context.Context, owner string, repo string, number int, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return c.clientREST.Issues.Edit(ctx, owner, repo, number, issue)
}

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

func (c *GithubClient) EditLabel(ctx context.Context, owner string, repo string, name string, label *github.Label) (*github.Label, *github.Response, error) {
	return c.clientREST.Issues.EditLabel(ctx, owner, repo, name, label)
}

func (c *GithubClient) RemoveLabelForIssue(ctx context.Context, owner string, repo string, number int, label string) (*github.Response, error) {
	return c.clientREST.Issues.RemoveLabelForIssue(ctx, owner, repo, number, label)
}

func (c *GithubClient) AddAssignees(ctx context.Context, owner string, repo string, number int, assignees []string) (*github.Issue, *github.Response, error) {
	return c.clientREST.Issues.AddAssignees(ctx, owner, repo, number, assignees)
}

func (c *GithubClient) ListIssuesByRepo(ctx context.Context, owner string, repo string, opts *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	is, err := PaginatedRequest(
		func() interface{} {
			return []*github.Issue{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			issues := i.([]*github.Issue)
			is, resp, err := c.clientREST.Issues.ListByRepo(ctx, owner, repo, opts)
			if err != nil {
				return nil, nil, err
			}
			issues = append(issues, is...)
			return issues, resp, nil
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return is.([]*github.Issue), nil, nil
}

func (c *GithubClient) GetComments(ctx context.Context, owner string, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, error) {
	fs, err := PaginatedRequest(
		func() interface{} {
			return []*github.IssueComment{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			fls := i.([]*github.IssueComment)
			fs, resp, err := c.clientREST.Issues.ListComments(ctx, owner, repo, number, &github.IssueListCommentsOptions{
				Sort:      opts.Sort,
				Direction: opts.Direction,
				Since:     opts.Since,
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: maxPerPage,
				},
			})
			if err != nil {
				return nil, nil, err
			}
			fls = append(fls, fs...)
			return fls, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return fs.([]*github.IssueComment), nil
}

func (c *GithubClient) GetLinkedProjectsForIssue(ctx context.Context, owner, repo string, number int, retryCount int) ([]GQLProjectV2Item, error) {
	projectItems := []GQLProjectV2Item{}
	hasNextPage := true
	currentRequestRetry := 1

	varGQLGetProjectFieldsQuery := map[string]interface{}{
		"repositoryOwner":    githubv4.String(owner),
		"repositoryName":     githubv4.String(repo),
		"issueNumber":        githubv4.Int(number),
		"projectItemsCursor": githubv4.String(""),
	}

	var getLinkedProjects IssueLinkedProjectsQuery

	for hasNextPage {
		if err := c.clientGQL.Query(ctx, &getLinkedProjects, varGQLGetProjectFieldsQuery); err != nil {
			currentRequestRetry++
			if currentRequestRetry <= retryCount {
				continue
			}
			return nil, err
		}

		items := getLinkedProjects.Repository.Issue.ProjectItems
		if items == nil {
			return nil, ErrProjectItemsNotFound
		}

		projectItems = append(projectItems, items.Nodes...)

		hasNextPage = items.PageInfo.HasNextPage

		varGQLGetProjectFieldsQuery["projectItemsCursor"] = githubv4.String(items.PageInfo.EndCursor)
	}

	return projectItems, nil
}
