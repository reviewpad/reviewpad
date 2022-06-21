// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/lang/aladino"
)

func mockClient(clientOptions ...mock.MockBackendOption) *http.Client {
	defaultMocks := []mock.MockBackendOption{
		mock.WithRequestMatchHandler(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				var pr int64 = 1234
				var label int64 = 208045946
				var url string = "https://api.github.com/repos/foobar/mocked-proj-1/pulls/1234"
				var date time.Time = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				w.Write(mock.MustMarshal(github.PullRequest{
					ID:   &pr,
					User: &github.User{Login: github.String("baz")},
					Assignees: []*github.User{
						{Login: github.String("foobar")},
					},
					Title:     github.String("Amazing new feature"),
					Body:      github.String("Please pull these awesome changes in!"),
					CreatedAt: &date,
					Commits:   github.Int(5),
					Number:    github.Int(1234),
					Milestone: &github.Milestone{
						Title: github.String("v1.0"),
					},
					Labels: []*github.Label{
						{
							ID: &label,
							Name: github.String("bug"),
							Description: github.String("Something isn't working"),
							Color: github.String("f29513"),
						},
					},
					Head: &github.PullRequestBranch{
						Repo: &github.Repository{
							Owner: &github.User{
								Login: github.String("foobar"),
							},
							URL: github.String(url),
							Name: github.String("mocked-proj-1"),
						},
						Ref: github.String("test"),
					},
					Base: &github.PullRequestBranch{
						Repo: &github.Repository{
							Owner: &github.User{
								Login: github.String("foobar"),
							},
							URL: github.String(url),
							Name: github.String("mocked-proj-1"),
						},
						Ref: github.String("main"),
					},
				}))
			}),
		),
		mock.WithRequestMatch(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			[]*github.CommitFile{
				{
					Filename: github.String("mocked-proj-1/1.txt"),
					Patch:    nil,
				},
				{
					Filename: github.String("mocked-proj-1/2.ts"),
					Patch:    nil,
				},
				{
					Filename: github.String("mocked-proj-1/docs/folder/README.md"),
					Patch:    nil,
				},
			},
		),
	}

	mocks := append(clientOptions, defaultMocks...)

	return mock.NewMockedHTTPClient(mocks...)
}

func mockEnv(clientOptions ...mock.MockBackendOption) (aladino.Env, error) {
	client := github.NewClient(mockClient(clientOptions...))

	owner := "foobar"
	repo := "mocked-proj-1"
	prNum := 6

	ctx := context.Background()

	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNum)
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
