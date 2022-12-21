// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var join = plugins_aladino.PluginBuiltIns().Functions["join"].Code

func TestJoin(t *testing.T) {
	tests := map[string]struct {
		args       []aladino.Value
		wantString *aladino.StringValue
	}{
		"join empty array": {
			args: []aladino.Value{
				aladino.BuildArrayValue([]aladino.Value{}),
				aladino.BuildStringValue(" "),
			},
			wantString: aladino.BuildStringValue(""),
		},
		"join single element array": {
			args: []aladino.Value{
				aladino.BuildArrayValue([]aladino.Value{
					aladino.BuildStringValue("a"),
				}),
				aladino.BuildStringValue(" "),
			},
			wantString: aladino.BuildStringValue("a"),
		},
		"join multiple elements array": {
			args: []aladino.Value{
				aladino.BuildArrayValue([]aladino.Value{
					aladino.BuildStringValue("a"),
					aladino.BuildStringValue("b"),
				}),
				aladino.BuildStringValue(" "),
			},
			wantString: aladino.BuildStringValue("a b"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)
			gotVal, err := join(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantString, gotVal)
		})
	}
}

func TestJoin_BadType(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)
	args := []aladino.Value{
		aladino.BuildArrayValue([]aladino.Value{
			aladino.BuildIntValue(4),
			aladino.BuildIntValue(2),
		}),
		aladino.BuildStringValue(" "),
	}
	got, err := join(mockedEnv, args)
	assert.Nil(t, got)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, fmt.Sprintf("join: invalid element of kind %v", aladino.INT_VALUE))
}
