// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var merge = plugins_aladino.PluginBuiltIns().Actions["merge"].Code

type MergeRequestPostBody struct {
	MergeMethod string `json:"merge_method"`
}

func TestMerge_WhenMergeMethodIsUnsupported(t *testing.T) {
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		State: github.String("open"),
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
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

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		State:  github.String("open"),
		Merged: github.Bool(false),
	})

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
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
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

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		State:  github.String("open"),
		Merged: github.Bool(false),
	})

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
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
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

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		State: github.String("open"),
		Draft: github.Bool(true),
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
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
}

func TestMerge_WhenMergeIsClosedPullRequest(t *testing.T) {
	wantMergeMethod := "rebase"

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		State: github.String("closed"),
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
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
}
