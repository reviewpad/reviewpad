// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var toNumber = plugins_aladino.PluginBuiltIns().Functions["toNumber"].Code

func TestToNumber(t *testing.T) {
	tests := map[string]struct {
		str     string
		wantRes aladino.Value
		wantErr error
	}{
		"when empty": {
			str: "",
			wantErr: fmt.Errorf(`error converting "" to number: %w`, &strconv.NumError{
				Func: "Atoi",
				Num:  "",
				Err:  errors.New("invalid syntax"),
			}),
		},
		"when not a number": {
			str: "a",
			wantErr: fmt.Errorf(`error converting "a" to number: %w`, &strconv.NumError{
				Func: "Atoi",
				Num:  "a",
				Err:  errors.New("invalid syntax"),
			}),
		},
		"when zero": {
			str:     "0",
			wantRes: aladino.BuildIntValue(0),
		},
		"when negative": {
			str:     "-1",
			wantRes: aladino.BuildIntValue(-1),
		},
		"when positive": {
			str:     "1",
			wantRes: aladino.BuildIntValue(1),
		},
		"when positive overflows": {
			str: "100000000000000000000000000000000000",
			wantErr: fmt.Errorf(`error converting "100000000000000000000000000000000000" to number: %w`, &strconv.NumError{
				Func: "Atoi",
				Num:  "100000000000000000000000000000000000",
				Err:  errors.New("value out of range"),
			}),
		},
		"when negative overflows": {
			str: "-100000000000000000000000000000000000",
			wantErr: fmt.Errorf(`error converting "-100000000000000000000000000000000000" to number: %w`, &strconv.NumError{
				Func: "Atoi",
				Num:  "-100000000000000000000000000000000000",
				Err:  errors.New("value out of range"),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

			res, err := toNumber(env, []aladino.Value{aladino.BuildStringValue(test.str)})

			assert.Equal(t, test.wantRes, res)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
