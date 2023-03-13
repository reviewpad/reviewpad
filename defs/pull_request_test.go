// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs_test

import (
	"testing"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/defs"
	"github.com/stretchr/testify/assert"
)

func TestFromGithubPullRequests(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		githubPullRequests []*github.PullRequest
		pullRequests       defs.PullRequests
	}{
		"when single pull request": {
			githubPullRequests: []*github.PullRequest{
				{
					NodeID: github.String("1"),
					Number: github.Int(1),
					State:  github.String("closed"),
					User: &github.User{
						NodeID: github.String("1"),
						Login:  github.String("login"),
					},
					Labels: []*github.Label{
						{
							NodeID: github.String("1"),
							Name:   github.String("label"),
						},
					},
					Body:  github.String("body"),
					Title: github.String("title"),
					RequestedReviewers: []*github.User{
						{
							NodeID: github.String("1"),
							Login:  github.String("login"),
						},
					},
					RequestedTeams: []*github.Team{
						{
							NodeID: github.String("1"),
							Name:   github.String("name"),
						},
					},
					Comments: github.Int(1),
					Commits:  github.Int(1),
					Base: &github.PullRequestBranch{
						Ref:   github.String("base"),
						SHA:   github.String("sha"),
						Label: github.String("label"),
						Repo: &github.Repository{
							NodeID: github.String("1"),
							Name:   github.String("name"),
							Owner: &github.User{
								NodeID: github.String("1"),
								Login:  github.String("login"),
							},
						},
					},
					Head: &github.PullRequestBranch{
						Ref:   github.String("head"),
						SHA:   github.String("sha"),
						Label: github.String("label"),
						Repo: &github.Repository{
							NodeID: github.String("1"),
							Name:   github.String("name"),
							Owner: &github.User{
								NodeID: github.String("1"),
								Login:  github.String("login"),
							},
						},
					},
					Draft: github.Bool(false),
					Assignees: []*github.User{
						{
							NodeID: github.String("1"),
							Login:  github.String("login"),
						},
					},
					CreatedAt:  &now,
					UpdatedAt:  &now,
					ClosedAt:   &now,
					Rebaseable: github.Bool(false),
					HTMLURL:    github.String("https://example.com"),
					Merged:     github.Bool(false),
					Milestone: &github.Milestone{
						NodeID: github.String("1"),
						Title:  github.String("title"),
					},
					Additions: github.Int(1),
					Deletions: github.Int(1),
				},
			},
			pullRequests: []defs.PullRequest{
				{
					ID:     "1",
					Number: 1,
					State:  "closed",
					User: defs.User{
						ID:    "1",
						Login: "login",
					},
					Labels: []defs.Label{
						{
							ID:   "1",
							Name: "label",
						},
					},
					Description: "body",
					Title:       "title",
					RequestedReviewers: []defs.User{
						{
							ID:    "1",
							Login: "login",
						},
					},
					RequestedTeams: []defs.Team{
						{
							ID:   "1",
							Name: "name",
						},
					},
					CommentsCount: 1,
					CommitsCount:  1,
					Repository: defs.Repository{
						ID:   "1",
						Name: "name",
						Owner: defs.User{
							ID:    "1",
							Login: "login",
						},
					},
					Assignees: defs.Users{
						defs.User{
							ID:    "1",
							Login: "login",
						},
					},
					Head: defs.Branch{
						Ref:   "head",
						SHA:   "sha",
						Label: "label",
						Repo: defs.Repository{
							ID:   "1",
							Name: "name",
							Owner: defs.User{
								ID:    "1",
								Login: "login",
							},
						},
					},
					Base: defs.Branch{
						Ref:   "base",
						SHA:   "sha",
						Label: "label",
						Repo: defs.Repository{
							ID:   "1",
							Name: "name",
							Owner: defs.User{
								ID:    "1",
								Login: "login",
							},
						},
					},
					CreatedAt:  now,
					UpdatedAt:  now,
					Closed:     true,
					ClosedAt:   &now,
					Rebaseable: false,
					URL:        "https://example.com",
					Merged:     false,
					Milestone: &defs.Milestone{
						ID:    "1",
						Title: "title",
					},
					Additions: 1,
					Deletions: 1,
				},
			},
		},
		"when multiple pull requests": {
			githubPullRequests: []*github.PullRequest{
				{
					NodeID: github.String("1"),
					Number: github.Int(1),
					State:  github.String("closed"),
					User: &github.User{
						NodeID: github.String("1"),
						Login:  github.String("login"),
					},
					Labels: []*github.Label{
						{
							NodeID: github.String("1"),
							Name:   github.String("label"),
						},
					},
					Body:  github.String("body"),
					Title: github.String("title"),
					RequestedReviewers: []*github.User{
						{
							NodeID: github.String("1"),
							Login:  github.String("login"),
						},
					},
					RequestedTeams: []*github.Team{
						{
							NodeID: github.String("1"),
							Name:   github.String("name"),
						},
					},
					Comments: github.Int(1),
					Commits:  github.Int(1),
					Base: &github.PullRequestBranch{
						Ref:   github.String("base"),
						SHA:   github.String("sha"),
						Label: github.String("label"),
						Repo: &github.Repository{
							NodeID: github.String("1"),
							Name:   github.String("name"),
							Owner: &github.User{
								NodeID: github.String("1"),
								Login:  github.String("login"),
							},
						},
					},
					Head: &github.PullRequestBranch{
						Ref:   github.String("head"),
						SHA:   github.String("sha"),
						Label: github.String("label"),
						Repo: &github.Repository{
							NodeID: github.String("1"),
							Name:   github.String("name"),
							Owner: &github.User{
								NodeID: github.String("1"),
								Login:  github.String("login"),
							},
						},
					},
					Draft: github.Bool(false),
					Assignees: []*github.User{
						{
							NodeID: github.String("1"),
							Login:  github.String("login"),
						},
					},
					CreatedAt: &now,
					UpdatedAt: &now,
					MergedAt:  &now,
					Merged:    github.Bool(true),
					ClosedAt:  &now,
				},
				{
					NodeID: github.String("2"),
					Number: github.Int(2),
					State:  github.String("open"),
					User: &github.User{
						NodeID: github.String("2"),
						Login:  github.String("login"),
					},
					Labels: []*github.Label{
						{
							NodeID: github.String("2"),
							Name:   github.String("label"),
						},
					},
					Body:  github.String("body"),
					Title: github.String("title"),
					RequestedReviewers: []*github.User{
						{
							NodeID: github.String("2"),
							Login:  github.String("login"),
						},
					},
					RequestedTeams: []*github.Team{},
					Comments:       github.Int(2),
					Commits:        github.Int(2),
					Base: &github.PullRequestBranch{
						Ref:   github.String("base"),
						SHA:   github.String("sha"),
						Label: github.String("label"),
						Repo: &github.Repository{
							NodeID: github.String("2"),
							Name:   github.String("name"),
							Owner: &github.User{
								NodeID: github.String("2"),
								Login:  github.String("login"),
							},
						},
					},
					Head: &github.PullRequestBranch{
						Ref:   github.String("head"),
						SHA:   github.String("sha"),
						Label: github.String("label"),
						Repo: &github.Repository{
							NodeID: github.String("2"),
							Name:   github.String("name"),
							Owner: &github.User{
								NodeID: github.String("2"),
								Login:  github.String("login"),
							},
						},
					},
					Draft: github.Bool(false),
					Assignees: []*github.User{
						{
							NodeID: github.String("2"),
							Login:  github.String("login"),
						},
					},
					CreatedAt: &now,
					UpdatedAt: &now,
				},
			},
			pullRequests: []defs.PullRequest{
				{
					ID:     "1",
					Number: 1,
					State:  "closed",
					User: defs.User{
						ID:    "1",
						Login: "login",
					},
					Labels: []defs.Label{
						{
							ID:   "1",
							Name: "label",
						},
					},
					Description: "body",
					Title:       "title",
					RequestedReviewers: []defs.User{
						{
							ID:    "1",
							Login: "login",
						},
					},
					RequestedTeams: []defs.Team{
						{
							ID:   "1",
							Name: "name",
						},
					},
					CommentsCount: 1,
					CommitsCount:  1,
					Repository: defs.Repository{
						ID:   "1",
						Name: "name",
						Owner: defs.User{
							ID:    "1",
							Login: "login",
						},
					},
					Assignees: defs.Users{
						defs.User{
							ID:    "1",
							Login: "login",
						},
					},
					Head: defs.Branch{
						Ref:   "head",
						SHA:   "sha",
						Label: "label",
						Repo: defs.Repository{
							ID:   "1",
							Name: "name",
							Owner: defs.User{
								ID:    "1",
								Login: "login",
							},
						},
					},
					Base: defs.Branch{
						Ref:   "base",
						SHA:   "sha",
						Label: "label",
						Repo: defs.Repository{
							ID:   "1",
							Name: "name",
							Owner: defs.User{
								ID:    "1",
								Login: "login",
							},
						},
					},
					CreatedAt: now,
					UpdatedAt: now,
					Merged:    true,
					MergedAt:  &now,
					Closed:    true,
					ClosedAt:  &now,
				},
				{
					ID:     "2",
					Number: 2,
					State:  "open",
					User: defs.User{
						ID:    "2",
						Login: "login",
					},
					Labels: []defs.Label{
						{
							ID:   "2",
							Name: "label",
						},
					},
					Description: "body",
					Title:       "title",
					RequestedReviewers: []defs.User{
						{
							ID:    "2",
							Login: "login",
						},
					},
					RequestedTeams: []defs.Team{},
					CommentsCount:  2,
					CommitsCount:   2,
					Repository: defs.Repository{
						ID:   "2",
						Name: "name",
						Owner: defs.User{
							ID:    "2",
							Login: "login",
						},
					},
					Assignees: defs.Users{
						defs.User{
							ID:    "2",
							Login: "login",
						},
					},
					Head: defs.Branch{
						Ref:   "head",
						SHA:   "sha",
						Label: "label",
						Repo: defs.Repository{
							ID:   "2",
							Name: "name",
							Owner: defs.User{
								ID:    "2",
								Login: "login",
							},
						},
					},
					Base: defs.Branch{
						Ref:   "base",
						SHA:   "sha",
						Label: "label",
						Repo: defs.Repository{
							ID:   "2",
							Name: "name",
							Owner: defs.User{
								ID:    "2",
								Login: "login",
							},
						},
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pullRequests := defs.FromGithubPullRequests(test.githubPullRequests)
			assert.Equal(t, test.pullRequests, pullRequests)
		})
	}
}
