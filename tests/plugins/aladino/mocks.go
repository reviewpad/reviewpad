// Copyright (C) 2019-2022 Explore.dev - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Copyright (C) 2019-2022 Explore.dev - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package aladino

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/engine"
)

func mockClient(clientOptions ...mock.MockBackendOption) *http.Client {
	defaultMocks := []mock.MockBackendOption{
		mock.WithRequestMatchHandler(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				var pr int64 = 1234
				var url string = "https://api.github.com/repos/foobar/mocked-proj-1/pulls/1234"
				var date time.Time = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
				w.Write(mock.MustMarshal(github.PullRequest{
					ID:   &pr,
					User: &github.User{Login: github.String("baz")},
					Assignees: []*github.User{
						{Login: github.String("foobar")},
					},
					Title:     github.String("A title"),
					Body:      github.String("A description"),
					CreatedAt: &date,
					Commits:   github.Int(5),
					Number:    github.Int(1234),
					Milestone: &github.Milestone{
						Title: github.String("A milestone"),
					},
					Labels: []*github.Label{
						{Name: github.String("default")},
					},
					Head: &github.PullRequestBranch{
						Repo: &github.Repository{
							Owner: &github.User{
								Login: github.String("foobar"),
							},
							URL: github.String(url),
						},
						Ref: github.String("test"),
					},
					Base: &github.PullRequestBranch{
						Repo: &github.Repository{
							Owner: &github.User{
								Login: github.String("foobar"),
							},
							URL: github.String(url),
						},
						Ref: github.String("main"),
					},
				}))
			}),
		),
	}

	mocks := append(clientOptions, defaultMocks...)

	return mock.NewMockedHTTPClient(mocks...)
}

func mockEnv(clientOptions ...mock.MockBackendOption) (*engine.Env, error) {
	client := github.NewClient(mockClient(clientOptions...))

	owner := "foobar"
	repo := "mocked-proj-1"
	prNum := 6

	ctx := context.Background()

	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNum)

	evalEnv := &engine.Env{
		Ctx:         ctx,
		Client:      client,
		ClientGQL:   nil,
		Collector:   nil,
		PullRequest: pr,
		Interpreter: nil,
	}

	return evalEnv, err
}