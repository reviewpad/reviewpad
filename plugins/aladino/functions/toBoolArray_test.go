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

var toBoolArray = plugins_aladino.PluginBuiltIns().Functions["toBoolArray"].Code

func TestToBoolArray(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes aladino.Value
		wantErr error
	}{
		"when empty": {
			str:     ``,
			wantErr: errors.New(`error converting "" to boolean array: unexpected end of JSON input`),
		},
		"when array of non boolean values": {
			str:     `[1, 2, "a", "b"]`,
			wantErr: errors.New(`error converting "[1, 2, "a", "b"]" to boolean array: json: cannot unmarshal number into Go value of type bool`),
		},
		"when array of boolean values": {
			str:     `[true, false, true]`,
			wantRes: aladino.BuildArrayValue([]aladino.Value{aladino.BuildBoolValue(true), aladino.BuildBoolValue(false), aladino.BuildBoolValue(true)}),
		},
		"when empty array": {
			str:     `[]`,
			wantRes: aladino.BuildArrayValue([]aladino.Value{}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toBoolArray(env, []aladino.Value{aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
