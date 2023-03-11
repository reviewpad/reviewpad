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
	RequestedTeams     Teams
	Assignees          Users
	CommentsCount      int
	CommitsCount       int
	Draft              bool
	Description        string
	Title              string
	Base               Branch
	Head               Branch
	Repository         Repository
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Merged             bool
	Closed             bool
	ClosedAt           *time.Time
	Rebaseable         bool
	URL                string
	MergedAt           *time.Time
	Milestone          *Milestone
	Additions          int
	Deletions          int
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
		RequestedTeams:     FromGithubTeams(githubPullRequest.RequestedTeams),
		CommentsCount:      githubPullRequest.GetComments(),
		CommitsCount:       githubPullRequest.GetCommits(),
		Repository:         FromGithubRepository(githubPullRequest.GetBase().GetRepo()),
		Draft:              githubPullRequest.GetDraft(),
		Assignees:          FromGithubUsers(githubPullRequest.Assignees),
		Base:               FromGithubBranch(githubPullRequest.GetBase()),
		Head:               FromGithubBranch(githubPullRequest.GetHead()),
		CreatedAt:          githubPullRequest.GetCreatedAt(),
		UpdatedAt:          githubPullRequest.GetUpdatedAt(),
		Merged:             githubPullRequest.GetMerged(),
		Closed:             githubPullRequest.GetState() == "closed",
		ClosedAt:           githubPullRequest.ClosedAt,
		Rebaseable:         githubPullRequest.GetRebaseable(),
		URL:                githubPullRequest.GetHTMLURL(),
		MergedAt:           githubPullRequest.MergedAt,
		Milestone:          FromGithubMilestone(githubPullRequest.Milestone),
		Additions:          githubPullRequest.GetAdditions(),
		Deletions:          githubPullRequest.GetDeletions(),
	}
}

func FromGithubPullRequests(githubPullRequests []*github.PullRequest) PullRequests {
	pullRequests := make(PullRequests, len(githubPullRequests))
	for i, githubPullRequest := range githubPullRequests {
		pullRequests[i] = FromGithubPullRequest(githubPullRequest)
	}
	return pullRequests
}
