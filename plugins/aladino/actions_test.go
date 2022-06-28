// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_test

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
	"github.com/reviewpad/reviewpad/v2/mocks"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/reviewpad/reviewpad/v2/utils/fmtio"
	"github.com/stretchr/testify/assert"
)

var addLabel = plugins_aladino.PluginBuiltIns().Actions["addLabel"].Code
var assignTeamReviewer = plugins_aladino.PluginBuiltIns().Actions["assignTeamReviewer"].Code
var commentOnce = plugins_aladino.PluginBuiltIns().Actions["commentOnce"].Code

type TeamReviewersRequestPostBody struct {
	TeamReviewers []string `json:"team_reviewers"`
}

func TestAssignTeamReviewer_WhenNoTeamSlugsAreProvided(t *testing.T) {
	mockedEnv, err := mocks.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})}
	err = assignTeamReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignTeamReviewer: requires at least 1 team to request for review")
}

func TestAssignTeamReviewer(t *testing.T) {
	teamA := "core"
	teamB := "reviewpad-project"
	wantTeamReviewers := []string{teamA, teamB}	
	gotTeamReviewers := []string{}
	mockedEnv, err := mocks.MockDefaultEnv(
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

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(teamA), aladino.BuildStringValue(teamB)})}
	err = assignTeamReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTeamReviewers, gotTeamReviewers)
}

func TestAddLabel_WhenANonStringArgIsProvided(t *testing.T) {
	mockedEnv, err := mocks.MockDefaultEnv()
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildIntValue(1)}
	err = addLabel(mockedEnv, args)

	assert.EqualError(t, err, "addLabel: expecting string argument, got IntValue")
}

func TestAddLabel_WhenGetLabelRequestFails(t *testing.T) {
	failMessage := "GetLabelRequestFail"
	mockedEnv, err := mocks.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetReposLabelsByOwnerByRepoByName,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildStringValue("test")}
	err = addLabel(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAddLabel_WhenAddLabelToIssueRequestFails(t *testing.T) {
	label := "bug"
	failMessage := "AddLabelsToIssueRequestFail"
	mockedEnv, err := mocks.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposLabelsByOwnerByRepoByName,
			&github.Label{
				Name: github.String(label),
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildStringValue(label)}
	err = addLabel(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAddLabel(t *testing.T) {
	labelA := "bug"
	wantLabels := []string{
		labelA,
	}
	gotLabels := []string{}
	mockedEnv, err := mocks.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposLabelsByOwnerByRepoByName,
			&github.Label{},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rawBody, _ := ioutil.ReadAll(r.Body)
				body := []string{}
				
				json.Unmarshal(rawBody, &body)

				gotLabels = body
			}),
		),
	)
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildStringValue(labelA)}
	err = addLabel(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabels, gotLabels)
}

func TestCommentOnce_WhenGetCommentsRequestFails(t *testing.T) {
	failMessage := "GetCommentRequestFail"
	comment := "Lorem Ipsum"
	mockedEnv, err := mocks.MockDefaultEnv(
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
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildStringValue(fmt.Sprintf("%v%v", plugins_aladino.ReviewpadCommentAnnotation, comment))}
	err = commentOnce(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestCommentOnce_WhenCommentAlreadyExists(t *testing.T) {
	existingComment := "Lorem Ipsum"
	commentCreated := false

	mockedEnv, err := mocks.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String(fmt.Sprintf("%v%v", plugins_aladino.ReviewpadCommentAnnotation, existingComment)),
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
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildStringValue(existingComment)}
	err = commentOnce(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, commentCreated, "The comment should not be created")
}

func TestCommentOnce_WhenFirstTime(t *testing.T) {
	commentToAdd := "Lorem Ipsum"
	addedComment := ""

	mockedEnv, err := mocks.MockDefaultEnv(
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
	fmtio.FailOnError("mockDefaultEnv failed: %v", err)

	args := []aladino.Value{aladino.BuildStringValue(commentToAdd)}
	err = commentOnce(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%v%v", plugins_aladino.ReviewpadCommentAnnotation, commentToAdd), addedComment)
}
