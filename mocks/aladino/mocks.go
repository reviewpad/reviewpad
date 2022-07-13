// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package mocks_aladino

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/shurcooL/githubv4"
)

const defaultMockPrID = 1234
const defaultMockPrNum = 6
const defaultMockPrOwner = "foobar"
const defaultMockPrRepoName = "default-mock-repo"

func GetDefaultMockPullRequestDetails() *github.PullRequest {
	prNum := defaultMockPrNum
	prId := int64(defaultMockPrID)
	prOwner := defaultMockPrOwner
	prRepoName := defaultMockPrRepoName
	prUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/%v", prOwner, prRepoName, prNum)
	prDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	return &github.PullRequest{
		ID:   &prId,
		User: &github.User{Login: github.String("john")},
		Assignees: []*github.User{
			{Login: github.String("jane")},
		},
		Title:     github.String("Amazing new feature"),
		Body:      github.String("Please pull these awesome changes in!"),
		CreatedAt: &prDate,
		Comments:  github.Int(6),
		Commits:   github.Int(5),
		Number:    github.Int(prNum),
		Milestone: &github.Milestone{
			Title: github.String("v1.0"),
		},
		Labels: []*github.Label{
			{
				Name: github.String("enhancement"),
			},
		},
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
		RequestedReviewers: []*github.User{
			{Login: github.String("jane")},
		},
	}
}

func GetDefaultMockPullRequestDetailsWith(pr *github.PullRequest) *github.PullRequest {
	defaultPullRequest := GetDefaultMockPullRequestDetails()

	if pr.User != nil {
		defaultPullRequest.User = pr.User
	}

	if pr.Number != nil {
		defaultPullRequest.Number = pr.Number
	}

	if pr.Base != nil {
		defaultPullRequest.Base = pr.Base
	}

	if pr.Assignees != nil {
		defaultPullRequest.Assignees = pr.Assignees
	}

    if pr.Commits != nil {
        defaultPullRequest.Commits = pr.Commits
    }
    
    if pr.RequestedReviewers != nil {
		defaultPullRequest.RequestedReviewers = pr.RequestedReviewers
	}

    if pr.RequestedReviewers != nil {
       defaultPullRequest.RequestedReviewers = pr.RequestedReviewers
    }

	return defaultPullRequest
}

func getDefaultMockPullRequestFileList() *[]*github.CommitFile {
	prRepoName := defaultMockPrRepoName
	return &[]*github.CommitFile{
		{
			Filename: github.String(fmt.Sprintf("%v/file1.ts", prRepoName)),
			Patch:    nil,
		},
		{
			Filename: github.String(fmt.Sprintf("%v/file2.ts", prRepoName)),
			Patch:    nil,
		},
		{
			Filename: github.String(fmt.Sprintf("%v/file3.ts", prRepoName)),
			Patch:    nil,
		},
	}
}

func mockHttpClientWith(clientOptions ...mock.MockBackendOption) *http.Client {
	return mock.NewMockedHTTPClient(clientOptions...)
}

func MockEnvWith(prOwner string, prRepoName string, prNum int, client *github.Client, clientGQL *githubv4.Client) (aladino.Env, error) {
	ctx := context.Background()
	pr, _, err := client.PullRequests.Get(ctx, prOwner, prRepoName, prNum)
	if err != nil {
		return nil, err
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		client,
		clientGQL,
		nil,
		pr,
		plugins_aladino.PluginBuiltIns(),
	)

	return env, err
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

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using handler directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

// MockDefaultEnv mocks an Aladino Env with default values.
func MockDefaultEnv(ghApiClientOptions []mock.MockBackendOption, ghGraphQLHandler func(http.ResponseWriter, *http.Request)) (aladino.Env, error) {
	prOwner := defaultMockPrOwner
	prRepoName := defaultMockPrRepoName
	prNum := defaultMockPrNum
	client := github.NewClient(mockDefaultHttpClient(ghApiClientOptions))

	// Handle GraphQL
    var clientGQL *githubv4.Client
    if ghGraphQLHandler != nil {
        mux := http.NewServeMux()
        mux.HandleFunc("/graphql", ghGraphQLHandler)
        clientGQL = githubv4.NewClient(&http.Client{Transport: localRoundTripper{handler: mux}})
    }

	return MockEnvWith(prOwner, prRepoName, prNum, client, clientGQL)
}

func MustRead(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}
