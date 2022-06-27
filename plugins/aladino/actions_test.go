// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

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
	"github.com/stretchr/testify/assert"
)

type TeamReviewersRequestPostBody struct {
	TeamReviewers []string `json:"team_reviewers"`
}

func TestAssignTeamReviewers_WhenArgIsNotArray(t *testing.T) {
	mockedEnv, err := mockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildIntValue(1)}
	err = assignTeamReviewerCode(mockedEnv, args)

	assert.EqualError(t, err, "assignTeamReviewer: requires array argument, got IntValue")
}

func TestAssignTeamReviewers(t *testing.T) {
	wantTeamReviewers := []string{
		"core",
		"reviewpad-project",
	}
	gotTeamReviewers := []string{}
	mockedEnv, err := mockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rawBody, _ := ioutil.ReadAll(r.Body)
				body := TeamReviewersRequestPostBody{}

				json.Unmarshal(rawBody, &body)

				gotTeamReviewers = body.TeamReviewers
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("core"), aladino.BuildStringValue("reviewpad-project")})}
	err = assignTeamReviewerCode(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTeamReviewers, gotTeamReviewers)
}

func TestCommentOnce_WhenGetCommentsRequestFails(t *testing.T) {
	failMessage := "GetCommentRequestFail"
	comment := "Lorem Ipsum"
	mockedEnv, err := mockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, comment))}
	err = commentOnceCode(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestCommentOnce_WhenCommentAlreadyExists(t *testing.T) {
	existingComment := "Lorem Ipsum"
	commentCreated := false

	mockedEnv, err := mockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(existingComment)}
	err = commentOnceCode(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, commentCreated, "The comment should not be created")
}

func TestCommentOnce_WhenFirstTime(t *testing.T) {
	commentToAdd := "Lorem Ipsum"
	addedComment := ""

	mockedEnv, err := mockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(commentToAdd)}
	err = commentOnceCode(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, commentToAdd), addedComment)
}
