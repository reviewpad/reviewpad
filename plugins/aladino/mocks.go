// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
)

const defaultMockPrNum = 123
const defaultMockPrOwner = "foobar"
const defaultMockPrRepoName = "default-mock-repo"

func getDefaultMockPullRequestDetails() *github.PullRequest {
	prNum := defaultMockPrNum
	prId := int64(prNum)
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
		Title:     github.String("Default Tile"),
		Body:      github.String("Default Body"),
		CreatedAt: &prDate,
		Commits:   github.Int(5),
		Number:    github.Int(prNum),
		Milestone: &github.Milestone{
			Title: github.String("v1.0"),
		},
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("headLogin"),
				},
				URL:  github.String(prUrl),
				Name: github.String(prRepoName),
			},
			Ref: github.String("HeadRef"),
		},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("baseLogin"),
				},
				URL:  github.String(prUrl),
				Name: github.String(prRepoName),
			},
			Ref: github.String("BaseRef"),
		},
	}
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

func mockEnvWith(prOwner string, prRepoName string, prNum int, client *github.Client) (aladino.Env, error) {
	ctx := context.Background()
	pr, _, err := client.PullRequests.Get(ctx, prOwner, prRepoName, prNum)
	if err != nil {
		return nil, err
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		client,
		nil,
		nil,
		pr,
		PluginBuiltIns(),
	)

	return env, err
}

// mockDefaultHttpClient mocks an HTTP client with default values and ready for Aladino.
// Being ready for Aladino means that at least the request to build Aladino Env need to be mocked.
// As for now, at least two github request need to be mocked in order to build Aladino Env, mainly:
// - Get pull request details (i.e. /repos/{owner}/{repo}/pulls/{pull_number})
// - Get pull request files (i.e. /repos/{owner}/{repo}/pulls/{pull_number})
func mockDefaultHttpClient(clientOptions ...mock.MockBackendOption) *http.Client {
	defaultMocks := []mock.MockBackendOption{
		mock.WithRequestMatchHandler(
			// Mock request to get pull request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(getDefaultMockPullRequestDetails()))
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

// mockDefaultEnv mocks an Aladino Env with default values.
func mockDefaultEnv(clientOptions ...mock.MockBackendOption) (aladino.Env, error) {
	prOwner := defaultMockPrOwner
	prRepoName := defaultMockPrRepoName
	prNum := defaultMockPrNum
	client := github.NewClient(mockDefaultHttpClient(clientOptions...))

	return mockEnvWith(prOwner, prRepoName, prNum, client)
}
