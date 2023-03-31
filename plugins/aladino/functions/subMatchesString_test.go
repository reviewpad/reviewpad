// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"regexp/syntax"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var subMatchesString = plugins_aladino.PluginBuiltIns().Functions["subMatchesString"].Code

func TestSubMatchesString(t *testing.T) {
	tests := map[string]struct {
		pattern    string
		str        string
		wantResult aladino.Value
		wantErr    error
	}{
		"when pattern is invalid": {
			pattern: "a(b",
			str:     "abc",
			wantErr: fmt.Errorf("failed to compile regex pattern a(b %w", &syntax.Error{
				Code: syntax.ErrorCode("missing closing )"),
				Expr: "a(b",
			}),
		},
		"when pattern is empty and string is empty": {
			pattern:    "",
			str:        "",
			wantResult: aladino.BuildArrayValue([]aladino.Value{}),
		},
		// empty regular expression matches everything
		"when pattern is empty but string has value": {
			pattern:    "",
			str:        "abc",
			wantResult: aladino.BuildArrayValue([]aladino.Value{}),
		},
		"when string is empty": {
			pattern:    "abc",
			str:        "",
			wantResult: aladino.BuildArrayValue([]aladino.Value{}),
		},
		"when pattern matches": {
			pattern:    "a(\\w+)bc(\\w+)",
			str:        "abcbcbc",
			wantResult: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("bc"), aladino.BuildStringValue("bc")}),
		},
		"when pattern doesn't match": {
			pattern:    "test123",
			str:        "[a-bA-Z]*",
			wantResult: aladino.BuildArrayValue([]aladino.Value{}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, nil, nil)

			res, err := subMatchesString(env, []aladino.Value{aladino.BuildStringValue(test.pattern), aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
