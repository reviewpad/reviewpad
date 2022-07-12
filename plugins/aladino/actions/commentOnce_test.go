// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

const ReviewpadCommentAnnotation = "<!--@annotation-reviewpad-single-comment-->"

var commentOnce = plugins_aladino.PluginBuiltIns().Actions["commentOnce"].Code

func TestCommentOnce_WhenGetCommentsRequestFails(t *testing.T) {
	failMessage := "GetCommentRequestFail"
	comment := "Lorem Ipsum"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, comment))}
	err = commentOnce(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestCommentOnce_WhenCommentAlreadyExists(t *testing.T) {
	existingComment := "Lorem Ipsum"
	commentCreated := false

	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						Body: github.String(fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, existingComment)),
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(existingComment)}
	err = commentOnce(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, commentCreated, "The comment should not be created")
}

func TestCommentOnce_WhenFirstTime(t *testing.T) {
	commentToAdd := "Lorem Ipsum"
	addedComment := ""

	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := github.IssueComment{}

					json.Unmarshal(rawBody, &body)

					addedComment = *body.Body
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(commentToAdd)}
	err = commentOnce(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, commentToAdd), addedComment)
}
