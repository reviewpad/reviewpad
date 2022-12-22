// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/ohler55/ojg/oj"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var toJSON = plugins_aladino.PluginBuiltIns().Functions["toJSON"].Code

func TestToJSON(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes aladino.Value
		wantErr error
	}{
		"when empty": {
			str:     ``,
			wantRes: aladino.BuildJSONValue(nil),
		},
		"when string": {
			str:     `"test"`,
			wantRes: aladino.BuildJSONValue("test"),
		},
		"when integer": {
			str:     `1`,
			wantRes: aladino.BuildJSONValue(int64(1)),
		},
		"when float": {
			str:     `1.0`,
			wantRes: aladino.BuildJSONValue(1.0),
		},
		"when true": {
			str:     `true`,
			wantRes: aladino.BuildJSONValue(true),
		},
		"when false": {
			str:     `false`,
			wantRes: aladino.BuildJSONValue(false),
		},
		"when array": {
			str:     `["a", "b", "c", true, false, [1, 2, 3]]`,
			wantRes: aladino.BuildJSONValue([]interface{}{"a", "b", "c", true, false, []interface{}{int64(1), int64(2), int64(3)}}),
		},
		"when object": {
			str: `{
				"id": 1,
				"name": "test",
				"isOwner": true,
				"labels": ["a", "b", "c"]
			}`,
			wantRes: aladino.BuildJSONValue(map[string]interface{}{
				"id":      int64(1),
				"name":    "test",
				"isOwner": true,
				"labels":  []interface{}{"a", "b", "c"},
			}),
		},
		"when invalid string": {
			str: `test`,
			wantErr: &oj.ParseError{
				Message: "expected true",
				Line:    1,
				Column:  2,
			},
		},
		"when invalid object syntax": {
			str: `{
				id: 1
			}`,
			wantErr: &oj.ParseError{
				Message: "expected a string start or object close, not 'i'",
				Line:    2,
				Column:  5,
			},
		},
		"when invalid array syntax": {
			str: `["1", "2"`,
			wantErr: &oj.ParseError{
				Message: "incomplete JSON",
				Line:    1,
				Column:  10,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toJSON(env, []aladino.Value{aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
