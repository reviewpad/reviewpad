// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
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

var sprintf = plugins_aladino.PluginBuiltIns().Functions["sprintf"].Code

func TestSprintf(t *testing.T) {
	tests := map[string]struct {
		format     string
		args       *lang.ArrayValue
		wantString *lang.StringValue
		wantErr    error
	}{
		"when in default format": {
			format: "Hello %v!",
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("world"),
			}),
			wantString: lang.BuildStringValue("Hello world!"),
		},
		"when in string and quotation format": {
			format: `Hello %q %s`,
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("world"),
				lang.BuildStringValue("again"),
			}),
			wantString: lang.BuildStringValue(`Hello "world" again`),
		},
		"when there are extra args": {
			format: "Hello %s",
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("earth"),
				lang.BuildStringValue("mars"),
				lang.BuildStringValue("venus"),
			}),
			wantString: lang.BuildStringValue("Hello earth"),
		},
		"when there are extra args with index formatting": {
			format: "Hello %[1]s",
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("earth"),
				lang.BuildStringValue("mars"),
				lang.BuildStringValue("venus"),
			}),
			wantString: lang.BuildStringValue("Hello earth"),
		},
		"when empty format and args": {
			format:     "",
			args:       lang.BuildArrayValue([]lang.Value{}),
			wantString: lang.BuildStringValue(""),
		},
		"when there is no formatting with args": {
			format: "Hello world",
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("earth"),
				lang.BuildStringValue("mars"),
				lang.BuildStringValue("venus"),
			}),
			wantString: lang.BuildStringValue("Hello world"),
		},
		"when formatting with numbers": {
			format: "Hello %s, You've added %d comments",
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("test"),
				lang.BuildIntValue(10),
			}),
			wantString: lang.BuildStringValue("Hello test, You've added 10 comments"),
		},
		"when formatting with numbers and booleans": {
			format: "hello %s, You've added %d comments, profile is locked: %t",
			args: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("test"),
				lang.BuildIntValue(10),
				lang.BuildBoolValue(true),
			}),
			wantString: lang.BuildStringValue("hello test, You've added 10 comments, profile is locked: true"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)
			gotString, err := sprintf(mockedEnv, []lang.Value{lang.BuildStringValue(test.format), test.args})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantString, gotString)
		})
	}
}
