// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
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

var selectFromJSON = plugins_aladino.PluginBuiltIns().Functions["selectFromJSON"].Code

func TestSelectFromJSON(t *testing.T) {
	exampleJSON := map[string]interface{}{
		"a": "b",
		"c": 1,
		"d": []interface{}{1, 2, 3},
		"e": true,
		"f": false,
		"g": 1.5,
		"nested": map[string]interface{}{
			"h": "test",
			"i": []interface{}{
				map[string]interface{}{
					"j": "k",
					"l": 10,
					"m": []interface{}{true, true, false, true, true},
				},
			},
		},
	}
	tests := map[string]struct {
		json       interface{}
		path       string
		wantResult aladino.Value
		wantErr    error
	}{
		"when nothing is found for null": {
			json:    nil,
			path:    "$.null",
			wantErr: errors.New(`nothing found at path "$.null"`),
		},
		"when nothing is found for string": {
			json:    "test",
			path:    "$.string",
			wantErr: errors.New(`nothing found at path "$.string"`),
		},
		"when nothing is found for int": {
			json:    1,
			path:    "$.number",
			wantErr: errors.New(`nothing found at path "$.number"`),
		},
		"when nothing is found for bool": {
			json:    true,
			path:    "$.bool",
			wantErr: errors.New(`nothing found at path "$.bool"`),
		},
		"when nothing is found for float": {
			json:    1.0,
			path:    "$.float",
			wantErr: errors.New(`nothing found at path "$.float"`),
		},
		"when nothing is found in array": {
			json:    []interface{}{"a", "b", "c", 1, 2, 3, true, false},
			path:    "$[10]",
			wantErr: errors.New(`nothing found at path "$[10]"`),
		},
		"when nothing is found in object": {
			json: map[string]interface{}{
				"a": "b",
				"c": 1,
				"d": []interface{}{1, 2, 3},
				"e": true,
				"f": false,
				"g": 1.0,
			},
			path:    "$.h",
			wantErr: errors.New(`nothing found at path "$.h"`),
		},
		"when string is found": {
			json:       exampleJSON,
			path:       "$.a",
			wantResult: aladino.BuildStringValue("b"),
		},
		"when int is found": {
			json:       exampleJSON,
			path:       "$.c",
			wantResult: aladino.BuildStringValue("1"),
		},
		"when bool is found": {
			json:       exampleJSON,
			path:       "$.e",
			wantResult: aladino.BuildStringValue("true"),
		},
		"when array is found": {
			json:       exampleJSON,
			path:       "$.d",
			wantResult: aladino.BuildStringValue("[1,2,3]"),
		},
		"when float is found": {
			json:       exampleJSON,
			path:       "$.g",
			wantResult: aladino.BuildStringValue("1.5"),
		},
		"when nested string is found with brace expression": {
			json:       exampleJSON,
			path:       `$["nested"]["h"]`,
			wantResult: aladino.BuildStringValue("test"),
		},
		"when nested array of objects is found with brace expression": {
			json:       exampleJSON,
			path:       `$["nested"]["i"]`,
			wantResult: aladino.BuildStringValue(`[{"j":"k","l":10,"m":[true,true,false,true,true]}]`),
		},
		"when string is found in nested object in array": {
			json:       exampleJSON,
			path:       `$.nested.i[0].j`,
			wantResult: aladino.BuildStringValue("k"),
		},
		"when true is found in nested object in array": {
			json:       exampleJSON,
			path:       `$.nested.i[0].m[0]`,
			wantResult: aladino.BuildStringValue("true"),
		},
		"when false is found in nested object in array": {
			json:       exampleJSON,
			path:       `$.nested.i[0].m[2]`,
			wantResult: aladino.BuildStringValue("false"),
		},
		"when root": {
			json:       exampleJSON,
			path:       `$`,
			wantResult: aladino.BuildStringValue(`{"a":"b","c":1,"d":[1,2,3],"e":true,"f":false,"g":1.5,"nested":{"h":"test","i":[{"j":"k","l":10,"m":[true,true,false,true,true]}]}}`),
		},
		"when root is null": {
			json:       nil,
			path:       `$`,
			wantResult: aladino.BuildStringValue("null"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, nil, nil)

			res, err := selectFromJSON(env, []aladino.Value{aladino.BuildJSONValue(test.json), aladino.BuildStringValue(test.path)})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
