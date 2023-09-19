// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/go-lib/entities"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewTypeEnv_WithDefaultEnv(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantTypeEnv := aladino.TypeEnv(map[string]lang.Type{
		"emptyFunction": lang.BuildFunctionType([]lang.Type{}, nil),
		"emptyAction":   lang.BuildFunctionType([]lang.Type{}, nil),
		"zeroConst":     lang.BuildFunctionType([]lang.Type{}, lang.BuildIntType()),
		"returnStr":     lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildStringType()),
	})

	gotTypeEnv := aladino.NewTypeEnv(mockedEnv)

	assert.Equal(t, wantTypeEnv, gotTypeEnv)
}

func TestNewEvalEnv_WhenGetPullRequestFilesFails(t *testing.T) {
	failMessage := "GetPullRequestFilesFail"
	mockedGithubClientREST := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			aladino.GetDefaultPullRequestDetails(),
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

	mockedGithubClient := gh.NewGithubClient(mockedGithubClientREST, nil, nil, nil)
	codehostClient := aladino.GetDefaultCodeHostClient(t, aladino.GetDefaultPullRequestDetails(), nil)
	ctx := context.Background()

	targetEntity := &entities.TargetEntity{
		Owner:  aladino.DefaultMockPrOwner,
		Repo:   aladino.DefaultMockPrRepoName,
		Number: aladino.DefaultMockPrNum,
		Kind:   entities.PullRequest,
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		mockedGithubClient,
		codehostClient,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
		nil,
		nil,
	)

	assert.Nil(t, env)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestNewEvalEnv_WhenNewFileFails(t *testing.T) {
	codehostClient := aladino.GetDefaultCodeHostClient(t, aladino.GetDefaultPullRequestDetails(), nil)

	githubClient := aladino.MockDefaultGithubClientWithFiles(nil, nil, []*github.CommitFile{
		{
			Patch: github.String("@@a"),
		},
	})

	ctx := context.Background()
	targetEntity := &entities.TargetEntity{
		Owner:  aladino.DefaultMockPrOwner,
		Repo:   aladino.DefaultMockPrRepoName,
		Number: aladino.DefaultMockPrNum,
		Kind:   entities.PullRequest,
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		githubClient,
		codehostClient,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
		nil,
		nil,
	)

	assert.Nil(t, env)
	assert.EqualError(t, err, "error in file patch : error in chunk lines parsing (1): missing lines info: @@a\npatch: @@a")
}

func TestNewEvalEnv(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	patch := "@@ -2,9 +2,11 @@ package main\n- func previous1() {\n+ func new1() {\n+\nreturn"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockedCodehostClient := aladino.GetDefaultCodeHostClient(t, aladino.GetDefaultPullRequestDetails(), nil)

	ctx := context.Background()
	mockedGithubClient := aladino.MockDefaultGithubClientWithFiles(nil, nil, []*github.CommitFile{
		{
			Filename: github.String(fileName),
			Patch:    github.String(patch),
		},
	})

	gotEnv, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		mockedGithubClient,
		mockedCodehostClient,
		aladino.DefaultMockCollector,
		aladino.DefaultMockTargetEntity,
		nil,
		aladino.MockBuiltIns(),
		nil,
		nil,
	)

	mockedTarget, _ := target.NewPullRequestTarget(ctx, aladino.DefaultMockTargetEntity, mockedGithubClient, mockedCodehostClient, mockedCodeReview)

	wantEnv := &aladino.BaseEnv{
		Ctx:          ctx,
		GithubClient: mockedGithubClient,
		Collector:    aladino.DefaultMockCollector,
		RegisterMap:  aladino.RegisterMap(make(map[string]lang.Value)),
		BuiltIns:     aladino.MockBuiltIns(),
		Report:       &aladino.Report{Actions: make([]string, 0)},
		// TODO: Mock an event
		EventPayload: nil,
		Target:       mockedTarget,
	}

	assert.Nil(t, err)
	assert.Equal(t, wantEnv.Ctx, gotEnv.GetCtx())
	assert.Equal(t, wantEnv.GithubClient, gotEnv.GetGithubClient())
	assert.Equal(t, wantEnv.Collector, gotEnv.GetCollector())
	assert.Equal(t, wantEnv.RegisterMap, gotEnv.GetRegisterMap())
	assert.Equal(t, wantEnv.EventPayload, gotEnv.GetEventPayload())
	// TODO: Compare the targets
	// assert.Equal(t, wantEnv.Target, gotEnv.GetTarget())

	assert.Equal(t, len(wantEnv.BuiltIns.Functions), len(gotEnv.GetBuiltIns().Functions))
	assert.Equal(t, len(wantEnv.BuiltIns.Actions), len(gotEnv.GetBuiltIns().Actions))
	for functionName, functionCode := range wantEnv.BuiltIns.Functions {
		assert.NotNil(t, gotEnv.GetBuiltIns().Functions[functionName])
		assert.Equal(t, functionCode.Type, gotEnv.GetBuiltIns().Functions[functionName].Type)
	}
	for actionName, actionCode := range wantEnv.BuiltIns.Actions {
		assert.NotNil(t, gotEnv.GetBuiltIns().Actions[actionName])
		assert.Equal(t, actionCode.Type, gotEnv.GetBuiltIns().Actions[actionName].Type)
	}

	assert.Equal(t, wantEnv.Report, gotEnv.GetReport())
}

func TestNewEvalEnv_WhenGetPullRequestFails(t *testing.T) {
	mockErr := errors.New("mock error")

	codehostClient := aladino.GetDefaultCodeHostClient(t, nil, mockErr)

	ctx := context.Background()

	targetEntity := &entities.TargetEntity{
		Owner:  aladino.DefaultMockPrOwner,
		Repo:   aladino.DefaultMockPrRepoName,
		Number: aladino.DefaultMockPrNum,
		Kind:   entities.PullRequest,
	}

	env, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		nil,
		codehostClient,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
		nil,
		nil,
	)

	assert.Nil(t, env)
	assert.Equal(t, err, mockErr)
}
