// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import (
	"time"

	"github.com/google/go-github/v49/github"
)

type PullRequest struct {
	ID                 string
	Number             int
	State              string
	User               User
	Labels             Labels
	RequestedReviewers Users
	Assignees          Users
	CommentsCount      int64
	CommitsCount       int64
	Draft              bool
	Description        string
	Title              string
	BaseBranch         string
	HeadBranch         string
	Repository         Repository
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type PullRequests []PullRequest

func FromGithubPullRequest(githubPullRequest *github.PullRequest) PullRequest {
	return PullRequest{
		ID:                 githubPullRequest.GetNodeID(),
		Number:             githubPullRequest.GetNumber(),
		State:              githubPullRequest.GetState(),
		User:               FromGithubUser(githubPullRequest.GetUser()),
		Labels:             FromGithubLabels(githubPullRequest.Labels),
		Description:        githubPullRequest.GetBody(),
		Title:              githubPullRequest.GetTitle(),
		RequestedReviewers: FromGithubUsers(githubPullRequest.RequestedReviewers),
		CommentsCount:      int64(githubPullRequest.GetComments()),
		CommitsCount:       int64(githubPullRequest.GetCommits()),
		Repository:         FromGithubRepository(githubPullRequest.GetBase().GetRepo()),
		Draft:              githubPullRequest.GetDraft(),
		Assignees:          FromGithubUsers(githubPullRequest.Assignees),
		BaseBranch:         githubPullRequest.GetBase().GetRef(),
		HeadBranch:         githubPullRequest.GetHead().GetRef(),
		CreatedAt:          githubPullRequest.GetCreatedAt(),
		UpdatedAt:          githubPullRequest.GetUpdatedAt(),
	}
}

func FromGithubPullRequests(githubPullRequests []*github.PullRequest) PullRequests {
	pullRequests := make(PullRequests, len(githubPullRequests))
	for i, githubPullRequest := range githubPullRequests {
		pullRequests[i] = FromGithubPullRequest(githubPullRequest)
	}
	return pullRequests
}
