// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"regexp/syntax"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var matchString = plugins_aladino.PluginBuiltIns().Functions["matchString"].Code

func TestMatchString(t *testing.T) {
	tests := map[string]struct {
		pattern    string
		str        string
		wantResult aladino.Value
		wantErr    error
	}{
		"when pattern is invalid": {
			pattern: "a(b",
			str:     "abc",
			wantErr: &syntax.Error{
				Code: syntax.ErrorCode("missing closing )"),
				Expr: "a(b",
			},
		},
		"when pattern is empty and string is empty": {
			pattern:    "",
			str:        "",
			wantResult: aladino.BuildBoolValue(true),
		},
		// empty regular expression matches everything
		"when pattern is empty but string has value": {
			pattern:    "",
			str:        "abc",
			wantResult: aladino.BuildBoolValue(true),
		},
		"when string is empty": {
			pattern:    "abc",
			str:        "",
			wantResult: aladino.BuildBoolValue(false),
		},
		"when pattern matches": {
			pattern:    "a(bc)+",
			str:        "abcbcbc",
			wantResult: aladino.BuildBoolValue(true),
		},
		"when pattern doesn't match": {
			pattern:    "test123",
			str:        "[a-bA-Z]*",
			wantResult: aladino.BuildBoolValue(false),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, nil, nil)

			res, err := matchString(env, []aladino.Value{aladino.BuildStringValue(test.pattern), aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
