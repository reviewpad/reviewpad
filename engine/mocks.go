// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/collector"
)

const defaultMockPrID = 1234
const defaultMockPrNum = 6
const defaultMockPrRepoName = "default-mock-repo"

var DefaultMockCtx = context.Background()
var DefaultMockCollector = collector.NewCollector("", "")
var MockedEventPayload = &github.CheckRunEvent{}

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

func getDefaultMockPullRequestFileList() *[]*github.CommitFile {
	prRepoName := defaultMockPrRepoName
	return &[]*github.CommitFile{
		{
			Filename: github.String(fmt.Sprintf("%v/file1.ts", prRepoName)),
			Patch:    github.String(`@@ -2,9 +2,11 @@ package main
- func previous1() {
+ func new1() {
+
return`),
		},
		{
			Filename: github.String(fmt.Sprintf("%v/file2.ts", prRepoName)),
			Patch:    github.String(`@@ -2,9 +2,11 @@ package main
- func previous2() {
+ func new2() {
+
return`),
		},
		{
			Filename: github.String(fmt.Sprintf("%v/file3.ts", prRepoName)),
			Patch:    github.String(`@@ -2,9 +2,11 @@ package main
- func previous3() {
+ func new3() {
+
return`),
		},
	}
}

// mockDefaultHttpClient mocks an HTTP client with default values and ready for Aladino.
// Being ready for Aladino means that at least the request to build Aladino Env need to be mocked.
// As for now, at least two github request need to be mocked in order to build Aladino Env, mainly:
// - Get pull request details (i.e. /repos/{owner}/{repo}/pulls/{pull_number})
// - Get pull request files (i.e. /repos/{owner}/{repo}/pulls/{pull_number}/files)
func mockDefaultHttpClient(clientOptions []mock.MockBackendOption) *http.Client {
	defaultMocks := []mock.MockBackendOption{
		mock.WithRequestMatchHandler(
			// Mock request to get pull request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(GetDefaultMockPullRequestDetails()))
			}),
		),
		mock.WithRequestMatchHandler(
			// Mock request to get pull request changed files
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(getDefaultMockPullRequestFileList()))
			}),
		),
	}

	mocks := append(clientOptions, defaultMocks...)

	return mockHttpClientWith(mocks...)
}

func mockHttpClientWith(clientOptions ...mock.MockBackendOption) *http.Client {
	return mock.NewMockedHTTPClient(clientOptions...)
}

func MockGithubClient(clientOptions []mock.MockBackendOption) *github.Client {
	return github.NewClient(mockDefaultHttpClient(clientOptions))
}
