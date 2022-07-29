// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var comments = plugins_aladino.PluginBuiltIns().Functions["comments"].Code

func TestComments(t *testing.T) {
	wantedComments := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("hello world"),
		},
	)

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						Body: github.String("hello world"),
					},
				},
			),
		},
		nil,
	)

	args := []aladino.Value{}
	gotComments, err := comments(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantedComments, gotComments)
}

func TestComments_WhenGetCommentsRequestFailed(t *testing.T) {
	failMessage := "GetCommentsRequestFailed"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
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
	)

	args := []aladino.Value{}
	gotComments, err := comments(mockedEnv, args)

	assert.Nil(t, gotComments)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}
