// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var titleLint = plugins_aladino.PluginBuiltIns().Actions["titleLint"].Code

func TestTitleLint(t *testing.T) {
	tests := map[string]struct {
		env                  aladino.Env
		wantReportedMessages map[aladino.Severity][]string
		wantErr              string
	}{
		"when type is empty": {
			env: aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
					Title: ": an amazing feature",
				}),
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{
				aladino.SEVERITY_ERROR: {
					"**Unconventional title detected**: ': an amazing feature' illegal ':' character in commit message type: col=00",
				},
			},
		},
		"when title is unconventional": {
			env: aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
					Title: "An Amazing Feature",
				}),
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{
				aladino.SEVERITY_ERROR: {
					"**Unconventional title detected**: 'An Amazing Feature' illegal 'A' character in commit message type: col=00",
				},
			},
		},
		"when title is conventional": {
			env: aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
					Title: "feat: an amazing feature",
				}),
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{},
		},
		"when title is conventional with scope": {
			env: aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				nil,
				nil,
				aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
					Title: "feat(api): an amazing feature",
				}),
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []lang.Value{}
			gotErr := titleLint(test.env, args)

			assert.Nil(t, gotErr)
			assert.Equal(t, test.wantReportedMessages, test.env.GetBuiltInsReportedMessages())
		})
	}
}
