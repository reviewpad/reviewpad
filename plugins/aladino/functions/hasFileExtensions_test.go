// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var hasFileExtensions = plugins_aladino.PluginBuiltIns().Functions["hasFileExtensions"].Code

func TestHasFileExtensions_WhenFalse(t *testing.T) {
	defaultMockPrFileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String(defaultMockPrFileName),
			Patch:    nil,
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(".md")})}
	gotVal, err := hasFileExtensions(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasFileExtensions_WhenTrue(t *testing.T) {
	defaultMockPrFileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := &[]*github.CommitFile{
		{
			Filename: github.String(defaultMockPrFileName),
			Patch:    nil,
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(".ts")})}
	gotVal, err := hasFileExtensions(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
