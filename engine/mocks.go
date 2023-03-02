// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/sirupsen/logrus"
)

// Use only for tests
const DefaultMockPrID = 1234
const DefaultMockPrNum = 6
const DefaultMockPrOwner = "foobar"
const DefaultMockPrRepoName = "default-mock-repo"
const DefaultMockEventName = "pull_request"
const DefaultMockEventAction = "opened"

// Use only for tests
var DefaultMockCtx = context.Background()
var DefaultMockLogger = logrus.NewEntry(logrus.New())
var DefaultMockCollector, _ = collector.NewCollector("", "distinctId", "pull_request", "runnerName", nil)
var DefaultMockEventPayload = &github.CheckRunEvent{}
var DefaultMockTargetEntity = &handler.TargetEntity{
	Owner:  DefaultMockPrOwner,
	Repo:   DefaultMockPrRepoName,
	Number: DefaultMockPrNum,
	Kind:   handler.PullRequest,
}
var DefaultMockEventDetails = &handler.EventDetails{
	EventName:   DefaultMockEventName,
	EventAction: DefaultMockEventAction,
}

func GetDefaultMockPullRequestDetails() *github.PullRequest {
	prNum := DefaultMockPrNum
	prId := int64(DefaultMockPrID)
	prOwner := DefaultMockPrOwner
	prRepoName := DefaultMockPrRepoName
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
			Filename: github.String(fmt.Sprintf("%v/file1.ts", DefaultMockPrRepoName)),
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
				MustWriteBytes(w, mock.MustMarshal(GetDefaultMockPullRequestDetails()))
			}),
		),
		mock.WithRequestMatchHandler(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				MustWriteBytes(w, mock.MustMarshal(getDefaultMockPullRequestFileList()))
			}),
		),
	}

	mocks := append(clientOptions, defaultMocks...)

	githubClientREST := github.NewClient(mock.NewMockedHTTPClient(mocks...))

	// TODO: mock the graphQL client
	return gh.NewGithubClient(githubClientREST, nil, nil)
}

func MockEnvWith(githubClient *gh.GithubClient, interpreter Interpreter, targetEntity *handler.TargetEntity, eventDetails *handler.EventDetails) (*Env, error) {
	dryRun := false
	mockedEnv, err := NewEvalEnv(
		DefaultMockCtx,
		DefaultMockLogger,
		dryRun,
		githubClient,
		DefaultMockCollector,
		targetEntity,
		interpreter,
		eventDetails,
	)

	if err != nil {
		return nil, fmt.Errorf("NewEvalEnv returned unexpected error: %v", err)
	}

	return mockedEnv, nil
}

func MustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func MustWriteBytes(w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		panic(err)
	}
}
