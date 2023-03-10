// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import (
	"time"

	"github.com/google/go-github/v49/github"
)

type Comment struct {
	ID        string
	Number    int64
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comments []Comment

func FromGithubComments(githubComments []*github.IssueComment) Comments {
	comments := make(Comments, len(githubComments))
	for i, githubComment := range githubComments {
		comments[i] = FromGithubComment(githubComment)
	}
	return comments
}

func FromGithubComment(githubComment *github.IssueComment) Comment {
	return Comment{
		ID:        githubComment.GetNodeID(),
		Number:    githubComment.GetID(),
		Body:      githubComment.GetBody(),
		CreatedAt: githubComment.GetCreatedAt(),
		UpdatedAt: githubComment.GetUpdatedAt(),
	}
}
