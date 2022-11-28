// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasBinaryFile = plugins_aladino.PluginBuiltIns().Functions["hasBinaryFile"].Code

func TestHasBinaryFile(t *testing.T) {
	tests := map[string]struct {
		wantResult *aladino.BoolValue
		wantErr    error
		files      []*github.CommitFile
	}{
		"when there is one file with patch": {
			wantResult: aladino.BuildBoolValue(false),
			files: []*github.CommitFile{
				{
					Filename: github.String("test.go"),
					Patch:    github.String("@@ -1,3 +1,3 @@ package test"),
				},
			},
		},
		"when there are two files with patch": {
			wantResult: aladino.BuildBoolValue(false),
			files: []*github.CommitFile{
				{
					Filename: github.String("test.go"),
					Patch:    github.String("@@ -1,3 +1,3 @@ package test"),
				},
				{
					Filename: github.String("test-2.go"),
					Patch:    github.String("@@ -1,3 +1,3 @@ package test_2"),
				},
			},
		},
		"when there is one file with no patch": {
			wantResult: aladino.BuildBoolValue(true),
			files: []*github.CommitFile{
				{
					Filename: github.String("test_bin_file"),
					Patch:    github.String(""),
				},
			},
		},
		"when there are two files one with no patch": {
			wantResult: aladino.BuildBoolValue(true),
			files: []*github.CommitFile{
				{
					Filename: github.String("test_bin_file"),
					Patch:    github.String(""),
				},
				{
					Filename: github.String("test-2.go"),
					Patch:    github.String("@@ -1,3 +1,3 @@ package test_2"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
					test.files,
				),
			}, nil, aladino.MockBuiltIns(), nil)

			res, err := hasBinaryFile(env, []aladino.Value{})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
