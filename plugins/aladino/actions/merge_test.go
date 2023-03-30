// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var merge = plugins_aladino.PluginBuiltIns().Actions["merge"].Code

type MergeRequestPostBody struct {
	MergeMethod string `json:"merge_method"`
}

func TestMerge_WhenMergeMethodIsUnsupported(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("INVALID")}
	err := merge(mockedEnv, args)

	assert.EqualError(t, err, "merge: unsupported merge method INVALID")
}

func TestMerge_WhenNoMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "merge"
	var gotMergeMethod string

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PutReposPullsMergeByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := MergeRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotMergeMethod = body.MergeMethod
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}

func TestMerge_WhenMergeMethodIsProvided(t *testing.T) {
	wantMergeMethod := "rebase"
	var gotMergeMethod string

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PutReposPullsMergeByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := MergeRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotMergeMethod = body.MergeMethod
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantMergeMethod, gotMergeMethod)
}

func TestMerge_WhenMergeIsOnDraftPullRequest(t *testing.T) {
	wantMergeMethod := "rebase"

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Status:  pbc.PullRequestStatus_OPEN,
		IsDraft: true,
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}

func TestMerge_WhenMergeIsClosedPullRequest(t *testing.T) {
	wantMergeMethod := "rebase"

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Status: pbc.PullRequestStatus_CLOSED,
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(wantMergeMethod)}
	err := merge(mockedEnv, args)

	assert.Nil(t, err)
}
