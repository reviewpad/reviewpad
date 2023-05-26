// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var join = plugins_aladino.PluginBuiltIns().Functions["join"].Code

func TestJoin(t *testing.T) {
	tests := map[string]struct {
		args       []lang.Value
		wantString *lang.StringValue
	}{
		"join empty array": {
			args: []lang.Value{
				lang.BuildArrayValue([]lang.Value{}),
				lang.BuildStringValue(" "),
			},
			wantString: lang.BuildStringValue(""),
		},
		"join single element array": {
			args: []lang.Value{
				lang.BuildArrayValue([]lang.Value{
					lang.BuildStringValue("a"),
				}),
				lang.BuildStringValue(" "),
			},
			wantString: lang.BuildStringValue("a"),
		},
		"join multiple elements array": {
			args: []lang.Value{
				lang.BuildArrayValue([]lang.Value{
					lang.BuildStringValue("a"),
					lang.BuildStringValue("b"),
				}),
				lang.BuildStringValue(" "),
			},
			wantString: lang.BuildStringValue("a b"),
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
	args := []lang.Value{
		lang.BuildArrayValue([]lang.Value{
			lang.BuildIntValue(4),
			lang.BuildIntValue(2),
		}),
		lang.BuildStringValue(" "),
	}
	got, err := join(mockedEnv, args)
	assert.Nil(t, got)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, fmt.Sprintf("join: invalid element of kind %v", lang.INT_VALUE))
}
