// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasLabel = plugins_aladino.PluginBuiltIns().Functions["hasLabel"].Code

func TestHasLabel(t *testing.T) {
	tests := map[string]struct {
		wantResult lang.Value
		wantErr    error
		label      string
	}{
		"when label exists": {
			wantResult: lang.BuildTrueValue(),
			label:      "large",
		},
		"when another label exists": {
			wantResult: lang.BuildTrueValue(),
			label:      "enhancement",
		},
		"when label does not exist": {
			wantResult: lang.BuildFalseValue(),
			label:      "small",
		},
	}
	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, []mock.MockBackendOption{}, nil, aladino.MockBuiltIns(), nil)

			res, err := hasLabel(env, []lang.Value{lang.BuildStringValue(test.label)})

			assert.Equal(t, test.wantResult, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
