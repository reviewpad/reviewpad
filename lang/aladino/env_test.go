// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"errors"
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
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
	mockError := errors.New(failMessage)
	codehostClient := aladino.GetDefaultCodeHostClientWithFiles(t, nil, mockError)
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
	assert.Equal(t, err.Error(), failMessage)
}

func TestNewEvalEnv_WhenNewFileFails(t *testing.T) {
	codehostClient := aladino.GetDefaultCodeHostClientWithFiles(t, []*pbc.File{
		{
			Patch: "@@a",
		},
	}, nil)

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
	assert.EqualError(t, err, "error in file patch : error in chunk lines parsing (1): missing lines info: @@a\npatch: @@a")
}

func TestNewEvalEnv(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	patch := "@@ -2,9 +2,11 @@ package main\n- func previous1() {\n+ func new1() {\n+\nreturn"

	mockedCodeReview := aladino.GetDefaultPullRequestDetails()

	mockedCodehostClient := aladino.GetDefaultCodeHostClientWithFiles(t, []*pbc.File{
		{
			Filename: fileName,
			Patch:    patch,
		},
	}, nil)

	ctx := context.Background()
	mockedGithubClient := gh.NewGithubClient(nil, nil, nil)

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

	codehostClient := aladino.GetDefaultCodeHostClient(t, nil, nil, mockErr, nil)

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
