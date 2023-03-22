// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
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
			mockedFiles := []*pbc.File{
				{
					Filename: "src/file.go",
				},
				{
					Filename: "test/main.go",
				},
				{
					Filename: "docs/main.md",
				},
				{
					Filename: "src/pkg/client.go",
				},
				{
					Filename: "src/pkg/client_test.go",
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
