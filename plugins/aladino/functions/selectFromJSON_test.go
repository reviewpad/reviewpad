// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
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
		wantResult lang.Value
		wantErr    error
	}{
		"when nothing is found in array": {
			json:       []interface{}{"a", "b", "c", 1, 2, 3, true, false},
			path:       "$[10]",
			wantResult: lang.BuildStringValue(""),
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
			path:       "$.h",
			wantResult: lang.BuildStringValue(""),
		},
		"when string is found": {
			json:       exampleJSON,
			path:       "$.a",
			wantResult: lang.BuildStringValue("b"),
		},
		"when int is found": {
			json:       exampleJSON,
			path:       "$.c",
			wantResult: lang.BuildStringValue("1"),
		},
		"when bool is found": {
			json:       exampleJSON,
			path:       "$.e",
			wantResult: lang.BuildStringValue("true"),
		},
		"when array is found": {
			json:       exampleJSON,
			path:       "$.d",
			wantResult: lang.BuildStringValue("[1,2,3]"),
		},
		"when float is found": {
			json:       exampleJSON,
			path:       "$.g",
			wantResult: lang.BuildStringValue("1.5"),
		},
		"when nested string is found with brace expression": {
			json:       exampleJSON,
			path:       `$["nested"]["h"]`,
			wantResult: lang.BuildStringValue("test"),
		},
		"when nested array of objects is found with brace expression": {
			json:       exampleJSON,
			path:       `$["nested"]["i"]`,
			wantResult: lang.BuildStringValue(`[{"j":"k","l":10,"m":[true,true,false,true,true]}]`),
		},
		"when string is found in nested object in array": {
			json:       exampleJSON,
			path:       `$.nested.i[0].j`,
			wantResult: lang.BuildStringValue("k"),
		},
		"when true is found in nested object in array": {
			json:       exampleJSON,
			path:       `$.nested.i[0].m[0]`,
			wantResult: lang.BuildStringValue("true"),
		},
		"when false is found in nested object in array": {
			json:       exampleJSON,
			path:       `$.nested.i[0].m[2]`,
			wantResult: lang.BuildStringValue("false"),
		},
		"when root": {
			json:       exampleJSON,
			path:       `$`,
			wantResult: lang.BuildStringValue(`{"a":"b","c":1,"d":[1,2,3],"e":true,"f":false,"g":1.5,"nested":{"h":"test","i":[{"j":"k","l":10,"m":[true,true,false,true,true]}]}}`),
		},
		"when root is null": {
			json:       nil,
			path:       `$`,
			wantResult: lang.BuildStringValue("null"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, nil, nil)

			res, err := selectFromJSON(env, []lang.Value{lang.BuildJSONValue(test.json), lang.BuildStringValue(test.path)})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
