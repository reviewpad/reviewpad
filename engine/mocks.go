// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/collector"
)

// Use only for tests
const defaultMockPrID = 1234
const defaultMockPrNum = 6

// Use only for tests
var DefaultMockCtx = context.Background()
var DefaultMockCollector = collector.NewCollector("", "")
var DefaultMockEventPayload = &github.CheckRunEvent{}

func GetDefaultMockPullRequestDetails() *github.PullRequest {
	prNum := defaultMockPrNum
	prId := int64(defaultMockPrID)
	prDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	return &github.PullRequest{
		ID:        &prId,
		User:      &github.User{Login: github.String("john")},
		Title:     github.String("Amazing new feature"),
		Body:      github.String("Please pull these awesome changes in!"),
		CreatedAt: &prDate,
		Number:    github.Int(prNum),
	}
}

func mockDefaultHttpClient(clientOptions []mock.MockBackendOption) *http.Client {
	return mockHttpClientWith(clientOptions...)
}

func mockHttpClientWith(clientOptions ...mock.MockBackendOption) *http.Client {
	return mock.NewMockedHTTPClient(clientOptions...)
}

func MockGithubClient(clientOptions []mock.MockBackendOption) *github.Client {
	return github.NewClient(mockDefaultHttpClient(clientOptions))
}
