// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	pbe "github.com/reviewpad/api/go/entities"
	api_mocks "github.com/reviewpad/api/go/mocks"
	pbs "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
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

func TestNewEvalEnv_WhenNewFileFails(t *testing.T) {
	hostsClient := api_mocks.NewMockHostsClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetCodeReview(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetCodeReviewReply{
			Review: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				Files: []*pbe.CommitFile{
					{
						Patch: "@@a",
					},
				},
			}),
		}, nil)

	codehostClient := &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}

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
		nil,
		codehostClient,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.EqualError(t, err, "error in file patch : error in chunk lines parsing (1): missing lines info: @@a\npatch: @@a")
}

func TestNewEvalEnv(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	patch := "@@ -2,9 +2,11 @@ package main\n- func previous1() {\n+ func new1() {\n+\nreturn"

	hostsClient := api_mocks.NewMockHostsClient(gomock.NewController(t))

	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Files: []*pbe.CommitFile{
			{
				Filename: fileName,
				Patch:    patch,
			},
		},
	})

	hostsClient.EXPECT().
		GetCodeReview(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&pbs.GetCodeReviewReply{
			Review: mockedCodeReview,
		}, nil)

	codehostClient := &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}

	ctx := context.Background()
	mockedGithubClient := gh.NewGithubClient(nil, nil, nil)

	gotEnv, err := aladino.NewEvalEnv(
		ctx,
		aladino.DefaultMockLogger,
		false,
		mockedGithubClient,
		codehostClient,
		aladino.DefaultMockCollector,
		aladino.DefaultMockTargetEntity,
		nil,
		aladino.MockBuiltIns(),
	)

	mockedTarget, _ := target.NewPullRequestTarget(ctx, aladino.DefaultMockTargetEntity, mockedGithubClient, mockedCodeReview)

	wantEnv := &aladino.BaseEnv{
		Ctx:          ctx,
		GithubClient: mockedGithubClient,
		Collector:    aladino.DefaultMockCollector,
		RegisterMap:  aladino.RegisterMap(make(map[string]aladino.Value)),
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
	hostsClient := api_mocks.NewMockHostsClient(gomock.NewController(t))

	hostsClient.EXPECT().
		GetCodeReview(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(nil, mockErr)

	codehostClient := &codehost.CodeHostClient{
		HostInfo: &codehost.HostInfo{
			Host:    pbe.Host_GITHUB,
			HostUri: "https://github.com",
		},
		CodehostClient: hostsClient,
	}

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
		nil,
		codehostClient,
		aladino.DefaultMockCollector,
		targetEntity,
		nil,
		aladino.MockBuiltIns(),
	)

	assert.Nil(t, env)
	assert.Equal(t, err, mockErr)
}
