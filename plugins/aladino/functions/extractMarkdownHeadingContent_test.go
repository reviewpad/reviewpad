// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var extractMarkdownHeadingContent = plugins_aladino.PluginBuiltIns().Functions["extractMarkdownHeadingContent"].Code

func TestExtractMarkdownHeadingContent(t *testing.T) {
	tests := map[string]struct {
		input        string
		headingTitle string
		headingLevel int
		wantResult   aladino.Value
		wantErr      error
	}{
		"when heading level is invalid": {
			input:        "a(b",
			headingTitle: "abc",
			headingLevel: 7,
			wantErr:      fmt.Errorf("invalid heading level: 7"),
		},
		"when input string is empty": {
			input:        "",
			headingTitle: "",
			headingLevel: 1,
			wantResult:   aladino.BuildStringValue(""),
		},
		"when heading 1 exists": {
			input:        "# Intro\nHello world\n\n## Intro\n\nSubsection\n",
			headingTitle: "Intro",
			headingLevel: 1,
			wantResult:   aladino.BuildStringValue("Hello world\n\n## Intro\n\nSubsection\n"),
		},
		"when heading 2 exists": {
			input:        "# Intro\nHello world\n\n## Intro one\n\nSubsection\n",
			headingTitle: "Intro one",
			headingLevel: 2,
			wantResult:   aladino.BuildStringValue("Subsection\n"),
		},
		"when heading 2 does not exist": {
			input:        "# Intro\nHello world\n\n## Intro one\n\nSubsection\n",
			headingTitle: "Intro",
			headingLevel: 2,
			wantResult:   aladino.BuildStringValue(""),
		},
		"system prompt exist": {
			input:        "# Intro\nHello world\n\n## System Prompt\n\nyou are awesome\n## User Prompt\n\nget chocolate\n",
			headingTitle: "System Prompt",
			headingLevel: 2,
			wantResult:   aladino.BuildStringValue("you are awesome"),
		},
		"user prompt exist": {
			input:        "# Intro\nHello world\n\n## System Prompt\n\nyou are awesome\n## User Prompt\n\nget chocolate\n## Test",
			headingTitle: "User Prompt",
			headingLevel: 2,
			wantResult:   aladino.BuildStringValue("get chocolate"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnv(t, nil, nil, nil, nil)

			res, err := extractMarkdownHeadingContent(env,
				[]aladino.Value{
					aladino.BuildStringValue(test.input),
					aladino.BuildStringValue(test.headingTitle),
					aladino.BuildIntValue(test.headingLevel),
				})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
