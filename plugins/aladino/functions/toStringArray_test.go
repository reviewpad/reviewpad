// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var toStringArray = plugins_aladino.PluginBuiltIns().Functions["toStringArray"].Code

func TestToStringArray(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes lang.Value
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
			wantRes: lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("a"), lang.BuildStringValue("b"), lang.BuildStringValue("c")}),
		},
		"when empty array": {
			str:     `[]`,
			wantRes: lang.BuildArrayValue([]lang.Value{}),
		},
		"when nested array": {
			str:     `[["a", "b", "c"]]`,
			wantErr: errors.New(`error converting "[["a", "b", "c"]]" to string array: json: cannot unmarshal array into Go value of type string`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toStringArray(env, []lang.Value{lang.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
