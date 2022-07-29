// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasFilePattern = plugins_aladino.PluginBuiltIns().Functions["hasFilePattern"].Code

func TestHasFilePattern_WhenFileBadPattern(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil)

	args := []aladino.Value{aladino.BuildStringValue("[0-9")}
	gotVal, err := hasFilePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Equal(t, wantVal, gotVal)
	assert.EqualError(t, err, "syntax error in pattern")
}

func TestHasFilePattern_WhenTrue(t *testing.T) {
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
					w.Write(mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("default-mock-repo/**")}
	gotVal, err := hasFilePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasFilePattern_WhenFalse(t *testing.T) {
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
					w.Write(mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
		},
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("default-mock-repo/test/**")}
	gotVal, err := hasFilePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
