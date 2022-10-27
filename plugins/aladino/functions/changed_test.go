// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var changed = plugins_aladino.PluginBuiltIns().Functions["changed"].Code

func TestChanged(t *testing.T) {
	tests := map[string]struct {
		args    []aladino.Value
		wantVal aladino.Value
	}{
		"bad spec": {
			args:    []aladino.Value{aladino.BuildStringValue("src/@1.go"), aladino.BuildStringValue("docs/file.md")},
			wantVal: aladino.BuildFalseValue(),
		},
		"missing in docs": {
			args:    []aladino.Value{aladino.BuildStringValue("src/@1.go"), aladino.BuildStringValue("docs/@1.md")},
			wantVal: aladino.BuildFalseValue(),
		},
		"changes in tests and docs": {
			args:    []aladino.Value{aladino.BuildStringValue("test/@1.go"), aladino.BuildStringValue("docs/@1.md")},
			wantVal: aladino.BuildTrueValue(),
		},
		"go tests": {
			args:    []aladino.Value{aladino.BuildStringValue("src/@1/@2.go"), aladino.BuildStringValue("src/@1/@2_test.go")},
			wantVal: aladino.BuildTrueValue(),
		},
		"nested patterns": {
			args:    []aladino.Value{aladino.BuildStringValue("src/pkg/@1.go"), aladino.BuildStringValue("src/pkg/dir/@2.go")},
			wantVal: aladino.BuildFalseValue(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedPullRequestFileList := []*github.CommitFile{
				{
					Filename: github.String("src/file.go"),
				},
				{
					Filename: github.String("test/main.go"),
				},
				{
					Filename: github.String("docs/main.md"),
				},
				{
					Filename: github.String("src/pkg/client.go"),
				},
				{
					Filename: github.String("src/pkg/client_test.go"),
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
				aladino.MockBuiltIns(),
				nil,
			)

			gotVal, err := changed(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
