// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var toBool = plugins_aladino.PluginBuiltIns().Functions["toBool"].Code

func TestToBool(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes aladino.Value
		wantErr error
	}{
		"when empty": {
			str: "",
			wantErr: fmt.Errorf(`error converting "" to boolean: %w`, &strconv.NumError{
				Func: "ParseBool",
				Num:  "",
				Err:  errors.New("invalid syntax"),
			}),
		},
		"when 1": {
			str:     "1",
			wantRes: aladino.BuildBoolValue(true),
		},
		"when 0": {
			str:     "0",
			wantRes: aladino.BuildBoolValue(false),
		},
		"when t": {
			str:     "t",
			wantRes: aladino.BuildBoolValue(true),
		},
		"when f": {
			str:     "f",
			wantRes: aladino.BuildBoolValue(false),
		},
		"when T": {
			str:     "T",
			wantRes: aladino.BuildBoolValue(true),
		},
		"when F": {
			str:     "F",
			wantRes: aladino.BuildBoolValue(false),
		},
		"when TRUE": {
			str:     "TRUE",
			wantRes: aladino.BuildBoolValue(true),
		},
		"when FALSE": {
			str:     "FALSE",
			wantRes: aladino.BuildBoolValue(false),
		},
		"when true": {
			str:     "true",
			wantRes: aladino.BuildBoolValue(true),
		},
		"when false": {
			str:     "false",
			wantRes: aladino.BuildBoolValue(false),
		},
		"when True": {
			str:     "True",
			wantRes: aladino.BuildBoolValue(true),
		},
		"when False": {
			str:     "False",
			wantRes: aladino.BuildBoolValue(false),
		},
		"when other": {
			str: "other",
			wantErr: fmt.Errorf(`error converting "other" to boolean: %w`, &strconv.NumError{
				Func: "ParseBool",
				Num:  "other",
				Err:  errors.New("invalid syntax"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toBool(env, []aladino.Value{aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
