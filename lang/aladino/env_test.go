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
