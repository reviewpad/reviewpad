// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewTypeEnv_WithDefaultEnv(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantTypeEnv := aladino.TypeEnv(map[string]aladino.Type{
		"emptyFunction": aladino.BuildFunctionType([]aladino.Type{}, nil),
		"emptyAction":   aladino.BuildFunctionType([]aladino.Type{}, nil),
		"zeroConst":     aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		"returnStr":     aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
	})

	gotTypeEnv := aladino.NewTypeEnv(mockedEnv)

	assert.Equal(t, wantTypeEnv, gotTypeEnv)
}

func TestNewEvalEnv_WhenGetPullRequestFilesFails(t *testing.T) {
	failMessage := "GetPullRequestFilesFail"
	mockedGithubClientREST := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			aladino.GetDefaultMockPullRequestDetails(),
		),
		mock.WithRequestMatchHandler(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	))

	mockedGithubClient := gh.NewGithubClient(mockedGithubClientREST, nil, nil)
	ctx := context.Background()

	targetEntity := &handler.TargetEntity{
		Owner:  aladino.DefaultMockPrOwner,
		Repo:   aladino.DefaultMockPrRepoName,
		Number: aladino.DefaultMockPrNum,
		Kind:   handler.PullRequest,
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		mockedGithubClient,
		nil,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestNewEvalEnv_WhenNewFileFails(t *testing.T) {
	mockedGithubClientREST := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			aladino.GetDefaultMockPullRequestDetails(),
		),
		mock.WithRequestMatch(
			mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
			&[]*github.CommitFile{
				{Patch: github.String("@@a")},
			},
		),
	))

	mockedGithubClient := gh.NewGithubClient(mockedGithubClientREST, nil, nil)
	ctx := context.Background()
	targetEntity := &handler.TargetEntity{
		Owner:  aladino.DefaultMockPrOwner,
		Repo:   aladino.DefaultMockPrRepoName,
		Number: aladino.DefaultMockPrNum,
		Kind:   handler.PullRequest,
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		mockedGithubClient,
		nil,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.EqualError(t, err, "error in file patch : error in chunk lines parsing (1): missing lines info: @@a\npatch: @@a")
}

func TestNewEvalEnv_WhenGetPullRequestFails(t *testing.T) {
	failMessage := "GetPullRequestFail"
	mockedGithubClientREST := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	))

	mockedGithubClient := gh.NewGithubClient(mockedGithubClientREST, nil, nil)
	ctx := context.Background()

	targetEntity := &handler.TargetEntity{
		Owner:  aladino.DefaultMockPrOwner,
		Repo:   aladino.DefaultMockPrRepoName,
		Number: aladino.DefaultMockPrNum,
		Kind:   handler.PullRequest,
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		mockedGithubClient,
		nil,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}
