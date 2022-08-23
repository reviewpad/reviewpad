// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/host-event-handler/handler"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
)

// Use only for tests
const defaultMockPrID = 1234
const defaultMockPrNum = 6
const defaultMockPrOwner = "foobar"
const defaultMockPrRepoName = "default-mock-repo"

// Use only for tests
var DefaultMockCtx = context.Background()
var DefaultMockCollector = collector.NewCollector("", "", "pull_request", "")
var DefaultMockEventPayload = &github.CheckRunEvent{}
var DefaultMockTargetEntity = &handler.TargetEntity{
	Owner:  defaultMockPrOwner,
	Repo:   defaultMockPrRepoName,
	Number: defaultMockPrNum,
	Kind:   handler.PullRequest,
}

func GetDefaultMockPullRequestDetails() *github.PullRequest {
	prNum := defaultMockPrNum
	prId := int64(defaultMockPrID)
	prOwner := defaultMockPrOwner
	prRepoName := defaultMockPrRepoName
	prDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	prUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/%v", prOwner, prRepoName, prNum)

	return &github.PullRequest{
		ID:        &prId,
		User:      &github.User{Login: github.String("john")},
		Title:     github.String("Amazing new feature"),
		Body:      github.String("Please pull these awesome changes in!"),
		CreatedAt: &prDate,
		Number:    github.Int(prNum),
		URL:       github.String(prUrl),
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("john"),
				},
				URL:  github.String(prUrl),
				Name: github.String(prRepoName),
			},
			Ref: github.String("new-topic"),
		},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("john"),
				},
				URL:  github.String(prUrl),
				Name: github.String(prRepoName),
			},
			Ref: github.String("master"),
		},
	}
}

func getDefaultMockPullRequestFileList() []*github.CommitFile {
	return []*github.CommitFile{
		{
			Filename: github.String(fmt.Sprintf("%v/file1.ts", defaultMockPrRepoName)),
			Patch: github.String(
				`@@ -2,9 +2,11 @@ package main
- func previous1() {
+ func new1() {
+
return`,
			),
		},
	}
}

func MockGithubClient(clientOptions []mock.MockBackendOption) *gh.GithubClient {
	defaultMocks := []mock.MockBackendOption{
		mock.WithRequestMatchHandler(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(GetDefaultMockPullRequestDetails()))
			}),
		),
		mock.WithRequestMatchHandler(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(getDefaultMockPullRequestFileList()))
			}),
		),
	}

	mocks := append(clientOptions, defaultMocks...)

	githubClientREST := github.NewClient(mock.NewMockedHTTPClient(mocks...))

	// TODO: mock the graphQL client
	return gh.NewGithubClient(githubClientREST, nil)
}

func MockEnvWith(githubClient *gh.GithubClient, interpreter Interpreter) (*Env, error) {
	dryRun := false
	mockedEnv, err := NewEvalEnv(
		DefaultMockCtx,
		dryRun,
		githubClient,
		DefaultMockCollector,
		DefaultMockTargetEntity,
		interpreter,
	)

	if err != nil {
		return nil, fmt.Errorf("NewEvalEnv returned unexpected error: %v", err)
	}

	return mockedEnv, nil
}
