// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var rebase = plugins_aladino.PluginBuiltIns().Actions["rebase"].Code

func TestRebase(t *testing.T) {
	prOwner := aladino.DefaultMockPrOwner
	prRepoName := aladino.DefaultMockPrRepoName
	prUrl := "https://github.com/reviewpad/TestGitRepository"

	tests := map[string]struct {
		mockedPullRequest *github.PullRequest
		wantErr           string
	}{
		"when pull request is not rebaseable": {
			mockedPullRequest: aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
				Rebaseable: github.Bool(false),
			}),
			wantErr: "the pull request is not rebaseable",
		},
		"when clone repository fails": {
			mockedPullRequest: aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
				Rebaseable: github.Bool(true),
			}),
			wantErr: "cannot set empty URL",
		},
		"when checkout branch fails": {
			mockedPullRequest: aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
				Rebaseable: github.Bool(true),
				Head: &github.PullRequestBranch{
					Repo: &github.Repository{
						Owner: &github.User{
							Login: github.String(prOwner),
						},
						HTMLURL: github.String(prUrl),
						Name:    github.String(prRepoName),
					},
					Ref: github.String("test"),
				},
			}),
			wantErr: "cannot locate remote-tracking branch 'origin/test'",
		},
		"when push fails": {
			mockedPullRequest: aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
				Rebaseable: github.Bool(true),
				Head: &github.PullRequestBranch{
					Repo: &github.Repository{
						Owner: &github.User{
							Login: github.String(prOwner),
						},
						HTMLURL: github.String(prUrl),
						Name:    github.String(prRepoName),
					},
					Ref: github.String("master"),
				},
			}),
			wantErr: "remote authentication required but no callback set",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Write(mock.MustMarshal(test.mockedPullRequest))
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			args := []aladino.Value{}
			gotErr := rebase(mockedEnv, args)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "rebase() error = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}
