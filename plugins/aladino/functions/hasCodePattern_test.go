// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var hasCodePattern = plugins_aladino.PluginBuiltIns().Functions["hasCodePattern"].Code

func TestHasCodePattern_WhenPullRequestPatchHasNilFile(t *testing.T) {
	fileName := "default-mock-repo/file1.ts"
	mockedPullRequestFileList := &[]*github.CommitFile{{
		Patch:    nil,
		Filename: github.String(fileName),
	}}
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

	mockedEnv.GetTarget().(*target.PullRequestTarget).Patch[fileName] = nil

	args := []aladino.Value{aladino.BuildStringValue("placeBet\\(.*\\)")}
	gotVal, err := hasCodePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasCodePattern_WhenPatternIsInvalid(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildStringValue("a(")}
	gotVal, err := hasCodePattern(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "query: compile error error parsing regexp: missing closing ): `a(`")
}

func TestHasCodePattern(t *testing.T) {
	mockedCodeReviewFileList := []*pbe.CommitFile{{
		Patch:    "@@ -2,9 +2,11 @@ package main\n- func previous() {\n+ func new() {\n+\nreturn",
		Filename: "default-mock-repo/file1.ts",
	}}
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
			Files: mockedCodeReviewFileList,
		}),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("new\\(.*\\)")}
	gotVal, err := hasCodePattern(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
