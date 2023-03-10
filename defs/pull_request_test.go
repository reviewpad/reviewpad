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
					State:  github.String("open"),
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
					Comments: github.Int(1),
					Commits:  github.Int(1),
					Base: &github.PullRequestBranch{
						Ref: github.String("base"),
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
						Ref: github.String("head"),
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
				},
			},
			pullRequests: []defs.PullRequest{
				{
					ID:     "1",
					Number: 1,
					State:  "open",
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
						{
							ID:    "1",
							Login: "login",
						},
					},
					BaseBranch: "base",
					HeadBranch: "head",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
		},
		"when multiple pull requests": {
			githubPullRequests: []*github.PullRequest{
				{
					NodeID: github.String("1"),
					Number: github.Int(1),
					State:  github.String("open"),
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
					Comments: github.Int(1),
					Commits:  github.Int(1),
					Base: &github.PullRequestBranch{
						Ref: github.String("base"),
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
						Ref: github.String("head"),
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
					Comments: github.Int(2),
					Commits:  github.Int(2),
					Base: &github.PullRequestBranch{
						Ref: github.String("base"),
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
						Ref: github.String("head"),
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
					State:  "open",
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
						{
							ID:    "1",
							Login: "login",
						},
					},
					BaseBranch: "base",
					HeadBranch: "head",
					CreatedAt:  now,
					UpdatedAt:  now,
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
					CommentsCount: 2,
					CommitsCount:  2,
					Repository: defs.Repository{
						ID:   "2",
						Name: "name",
						Owner: defs.User{
							ID:    "2",
							Login: "login",
						},
					},
					Assignees: defs.Users{
						{
							ID:    "2",
							Login: "login",
						},
					},
					BaseBranch: "base",
					HeadBranch: "head",
					CreatedAt:  now,
					UpdatedAt:  now,
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
