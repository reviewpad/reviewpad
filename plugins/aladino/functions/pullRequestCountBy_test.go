// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var pullRequestCountBy = plugins_aladino.PluginBuiltIns().Functions["pullRequestCountBy"].Code

func TestPullRequestCountBy_WhenListIssuesByRepoFails(t *testing.T) {
	failMessage := "ListListIssuesByRepoRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposIssuesByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{lang.BuildStringValue(""), lang.BuildStringValue("")}
	gotTotal, err := pullRequestCountBy(mockedEnv, args)

	assert.Nil(t, gotTotal)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestPullRequestCountBy(t *testing.T) {
	firstIssue := &github.Issue{
		Title: github.String("First Issue"),
		State: github.String("open"),
	}

	secondIssue := &github.Issue{
		Title: github.String("Second Issue"),
		PullRequestLinks: &github.PullRequestLinks{
			URL: github.String("pull-request-link"),
		},
		State: github.String("closed"),
	}

	thirdIssue := &github.Issue{
		Title: github.String("Third Issue"),
		State: github.String("closed"),
		User: &github.User{
			Login: github.String("steve"),
		},
	}

	tests := map[string]struct {
		args    []lang.Value
		issues  []*github.Issue
		wantVal lang.Value
	}{
		"default values": {
			args:    []lang.Value{lang.BuildStringValue(""), lang.BuildStringValue("")},
			issues:  []*github.Issue{firstIssue, secondIssue, thirdIssue},
			wantVal: lang.BuildIntValue(1),
		},
		"only closed pull requests": {
			args:    []lang.Value{lang.BuildStringValue(""), lang.BuildStringValue("closed")},
			issues:  []*github.Issue{firstIssue, secondIssue},
			wantVal: lang.BuildIntValue(1),
		},
		"only pull requests by steve": {
			args:    []lang.Value{lang.BuildStringValue("steve"), lang.BuildStringValue("")},
			issues:  []*github.Issue{firstIssue, thirdIssue},
			wantVal: lang.BuildIntValue(0),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposIssuesByOwnerByRepo,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(test.issues))
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			gotVal, err := pullRequestCountBy(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
