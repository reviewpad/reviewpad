// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var isDraft = plugins_aladino.PluginBuiltIns().Functions["isDraft"].Code

func TestIsDraftOnNilPullRequest(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedEnvWithNilPullRequest := &aladino.BaseEnv{
		Ctx:         mockedEnv.GetCtx(),
		Client:      mockedEnv.GetClient(),
		ClientGQL:   mockedEnv.GetClientGQL(),
		Collector:   mockedEnv.GetCollector(),
		PullRequest: nil,
		Patch:       mockedEnv.GetPatch(),
		RegisterMap: mockedEnv.GetRegisterMap(),
		BuiltIns:    mockedEnv.GetBuiltIns(),
		Report:      mockedEnv.GetReport(),
	}

	args := []aladino.Value{}
	gotVal, err := isDraft(mockedEnvWithNilPullRequest, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "isDraft: pull request is nil")
}

func TestIsDraft_WhenTrue(t *testing.T) {
    mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    mockedEnv.GetPullRequest().Draft = github.Bool(true)

	args := []aladino.Value{}
	gotVal, err := isDraft(mockedEnv, args)

    wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal, "pull request should be in draft")
}

func TestIsDraft_WhenFalse(t *testing.T) {
    mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    mockedEnv.GetPullRequest().Draft = github.Bool(false)

	args := []aladino.Value{}
	gotVal, err := isDraft(mockedEnv, args)

    wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal, "pull request should not be in draft")
}

