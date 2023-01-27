// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

const ReviewpadCommentAnnotation = "<!--@annotation-reviewpad-single-comment-->"

var commentOnce = plugins_aladino.PluginBuiltIns().Actions["commentOnce"].Code

func TestCommentOnce_WhenGetCommentsRequestFails(t *testing.T) {
	failMessage := "GetCommentRequestFail"
	comment := "Lorem Ipsum"
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
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(fmt.Sprintf("%s%s", ReviewpadCommentAnnotation, comment))}
	err := commentOnce(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestCommentOnce_WhenCommentAlreadyExists(t *testing.T) {
	existingComment := "Lorem Ipsum"
	commentCreated := false

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						Body: github.String(fmt.Sprintf("%s%s\n%s", ReviewpadCommentAnnotation, aladino.ReviewpadIgnoreCommentAnnotation, existingComment)),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the create comment request was performed then the comment was created
					commentCreated = true
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(existingComment)}
	err := commentOnce(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, commentCreated, "The comment should not be created")
}

func TestCommentOnce_WhenFirstTime(t *testing.T) {
	commentToAdd := "Lorem Ipsum"
	addedComment := ""

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := github.IssueComment{}

					utils.MustUnmarshal(rawBody, &body)

					addedComment = *body.Body
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue(commentToAdd)}
	err := commentOnce(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s%s\n%s", ReviewpadCommentAnnotation, aladino.ReviewpadIgnoreCommentAnnotation, commentToAdd), addedComment)
}
