// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var changed = plugins_aladino.PluginBuiltIns().Functions["changed"].Code

func TestChanged(t *testing.T) {
	tests := map[string]struct {
		args    []lang.Value
		wantVal lang.Value
	}{
		"bad spec": {
			args:    []lang.Value{lang.BuildStringValue("src/@1.go"), lang.BuildStringValue("docs/file.md")},
			wantVal: lang.BuildFalseValue(),
		},
		"missing in docs": {
			args:    []lang.Value{lang.BuildStringValue("src/@1.go"), lang.BuildStringValue("docs/@1.md")},
			wantVal: lang.BuildFalseValue(),
		},
		"changes in tests and docs": {
			args:    []lang.Value{lang.BuildStringValue("test/@1.go"), lang.BuildStringValue("docs/@1.md")},
			wantVal: lang.BuildTrueValue(),
		},
		"go tests": {
			args:    []lang.Value{lang.BuildStringValue("src/@1/@2.go"), lang.BuildStringValue("src/@1/@2_test.go")},
			wantVal: lang.BuildTrueValue(),
		},
		"nested patterns": {
			args:    []lang.Value{lang.BuildStringValue("src/pkg/@1.go"), lang.BuildStringValue("src/pkg/dir/@2.go")},
			wantVal: lang.BuildFalseValue(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedFiles := []*github.CommitFile{
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

			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				aladino.GetDefaultPullRequestDetails(),
				mockedFiles,
				aladino.MockBuiltIns(),
				nil,
			)

			gotVal, err := changed(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
