// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var toStringArray = plugins_aladino.PluginBuiltIns().Functions["toStringArray"].Code

func TestToStringArray(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes aladino.Value
		wantErr error
	}{
		"when empty": {
			str:     ``,
			wantErr: errors.New(`error converting "" to string array: unexpected end of JSON input`),
		},
		"when array of non string values": {
			str:     `[1, 2, true]`,
			wantErr: errors.New(`error converting "[1, 2, true]" to string array: json: cannot unmarshal number into Go value of type string`),
		},
		"when array of strings": {
			str:     `["a", "b", "c"]`,
			wantRes: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("a"), aladino.BuildStringValue("b"), aladino.BuildStringValue("c")}),
		},
		"when empty array": {
			str:     `[]`,
			wantRes: aladino.BuildArrayValue([]aladino.Value{}),
		},
		"when nested array": {
			str:     `[["a", "b", "c"]]`,
			wantErr: errors.New(`error converting "[["a", "b", "c"]]" to string array: json: cannot unmarshal array into Go value of type string`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toStringArray(env, []aladino.Value{aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
