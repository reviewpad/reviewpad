// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewInterpreter_WhenNewEvalEnvFails(t *testing.T) {
	ctx := context.Background()
	failMessage := "GetPullRequestFilesRequestFail"
	client := github.NewClient(
		mock.NewMockedHTTPClient(
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
		),
	)
	gotInterpreter, err := aladino.NewInterpreter(
		ctx,
		client,
		nil,
		nil,
		mocks_aladino.GetDefaultMockPullRequestDetails(),
		nil,
	)

	assert.Nil(t, gotInterpreter)
	assert.NotNil(t, err)
}

func TestNewInterpreter(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	gotInterpreter, err := aladino.NewInterpreter(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedEnv.GetClientGQL(),
		mockedEnv.GetCollector(),
		mockedEnv.GetPullRequest(),
		mockedEnv.GetBuiltIns(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantInterpreter, gotInterpreter)
}
