// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var commitLint = plugins_aladino.PluginBuiltIns().Actions["commitLint"].Code

func TestCommitLint_WhenRequestFails(t *testing.T) {
	tests := map[string]struct {
		env                  aladino.Env
		wantReportedMessages map[aladino.Severity][]string
		wantErr              string
	}{
		"when get pull request commits fails": {
			env: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"GetPullRequestCommitsRequestFailed",
							)
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{},
			wantErr:              "GetPullRequestCommitsRequestFailed",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{}
			gotErr := commitLint(test.env, args)

			gotReportedMessages := test.env.GetBuiltInsReportedMessages()

			assert.Equal(t, test.wantReportedMessages, gotReportedMessages)
			assert.Equal(t, gotErr.(*github.ErrorResponse).Message, test.wantErr)
		})
	}
}

func TestCommitLint(t *testing.T) {
	tests := map[string]struct {
		env                  aladino.Env
		wantReportedMessages map[aladino.Severity][]string
		wantErr              string
	}{
		"when commit parser fails": {
			env: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						[]*github.RepositoryCommit{
							{
								SHA: github.String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
								// illegal '&' character in commit message type
								Commit: &github.Commit{Message: github.String("&")},
							},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{
				aladino.SEVERITY_ERROR: {
					"**Unconventional commit detected**: '&' (6dcb09b5b57875f334f61aebed695e2e4193db5e)",
				},
			},
		},
		"when there is one non conventional commit": {
			env: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						[]*github.RepositoryCommit{
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
								Commit: &github.Commit{Message: github.String("Fix all bugs")},
							},
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5a"),
								Commit: &github.Commit{Message: github.String("docs: correct spelling of CHANGELOG")},
							},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{
				aladino.SEVERITY_ERROR: {
					"**Unconventional commit detected**: 'Fix all bugs' (6dcb09b5b57875f334f61aebed695e2e4193db5e)",
				},
			},
		},
		"when there is multiple non conventional commits": {
			env: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						[]*github.RepositoryCommit{
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
								Commit: &github.Commit{Message: github.String("Fix all bugs")},
							},
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5a"),
								Commit: &github.Commit{Message: github.String("docs: correct spelling of CHANGELOG")},
							},
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5c"),
								Commit: &github.Commit{Message: github.String("Add code comments")},
							},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{
				aladino.SEVERITY_ERROR: {
					"**Unconventional commit detected**: 'Fix all bugs' (6dcb09b5b57875f334f61aebed695e2e4193db5e)",
					"**Unconventional commit detected**: 'Add code comments' (6dcb09b5b57875f334f61aebed695e2e4193db5c)",
				},
			},
		},
		"when there is no non conventional commits": {
			env: aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						[]*github.RepositoryCommit{
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5a"),
								Commit: &github.Commit{Message: github.String("docs: correct spelling of CHANGELOG")},
							},
							{
								SHA:    github.String("6dcb09b5b57875f334f61aebed695e2e4193db5c"),
								Commit: &github.Commit{Message: github.String("feat(lang): add Polish language")},
							},
						},
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			),
			wantReportedMessages: map[aladino.Severity][]string{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			args := []aladino.Value{}
			gotErr := commitLint(test.env, args)

			assert.Nil(t, gotErr)
			assert.Equal(t, test.wantReportedMessages, test.env.GetBuiltInsReportedMessages())
		})
	}
}
