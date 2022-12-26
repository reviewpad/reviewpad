// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var toIntArray = plugins_aladino.PluginBuiltIns().Functions["toIntArray"].Code

func TestToIntArray(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes aladino.Value
		wantErr error
	}{
		"when empty": {
			str:     ``,
			wantErr: errors.New(`error converting "" to int array: unexpected end of JSON input`),
		},
		"when array of non string values": {
			str:     `["a", "b", true]`,
			wantErr: errors.New(`error converting "["a", "b", true]" to int array: json: cannot unmarshal string into Go value of type int`),
		},
		"when array of ints": {
			str:     `[1, 2, 3]`,
			wantRes: aladino.BuildArrayValue([]aladino.Value{aladino.BuildIntValue(1), aladino.BuildIntValue(2), aladino.BuildIntValue(3)}),
		},
		"when empty array": {
			str:     `[]`,
			wantRes: aladino.BuildArrayValue([]aladino.Value{}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toIntArray(env, []aladino.Value{aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
