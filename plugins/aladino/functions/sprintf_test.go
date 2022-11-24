// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var sprintf = plugins_aladino.PluginBuiltIns().Functions["sprintf"].Code

func TestSprintf(t *testing.T) {
	tests := map[string]struct {
		format     string
		args       *aladino.ArrayValue
		wantString *aladino.StringValue
		wantErr    error
	}{
		"when in default format": {
			format: "Hello %v!",
			args: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("world"),
			}),
			wantString: aladino.BuildStringValue("Hello world!"),
		},
		"when in string and quotation format": {
			format: `Hello %q %s`,
			args: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("world"),
				aladino.BuildStringValue("again"),
			}),
			wantString: aladino.BuildStringValue(`Hello "world" again`),
		},
		"when there are extra args": {
			format: "Hello %s",
			args: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("earth"),
				aladino.BuildStringValue("mars"),
				aladino.BuildStringValue("venus"),
			}),
			wantString: aladino.BuildStringValue("Hello earth"),
		},
		"when there are extra args with index formatting": {
			format: "Hello %[1]s",
			args: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("earth"),
				aladino.BuildStringValue("mars"),
				aladino.BuildStringValue("venus"),
			}),
			wantString: aladino.BuildStringValue("Hello earth"),
		},
		"when empty format and args": {
			format:     "",
			args:       aladino.BuildArrayValue([]aladino.Value{}),
			wantString: aladino.BuildStringValue(""),
		},
		"when there is no formatting with args": {
			format: "Hello world",
			args: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("earth"),
				aladino.BuildStringValue("mars"),
				aladino.BuildStringValue("venus"),
			}),
			wantString: aladino.BuildStringValue("Hello world"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)
			gotString, err := sprintf(mockedEnv, []aladino.Value{aladino.BuildStringValue(test.format), test.args})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantString, gotString)
		})
	}
}
