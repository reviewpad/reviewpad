// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v45/github"
	"github.com/gorilla/mux"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
)

type EditCommentRequestPostBody struct {
	Body string `json:"body"`
}

func TestReportError(t *testing.T) {
	errorMsg := "Testing Error"

	wantErr := fmt.Sprintf("[report] %v", errorMsg)
	gotErr := reportError(errorMsg)

	assert.EqualError(t, gotErr, wantErr)
}

func TestAddToReport(t *testing.T) {
	statement := engine.BuildStatement("$addLabel(\"test\")")

	report := Report{
		Actions: []string{statement.GetStatementCode()},
	}

	wantReport := Report{
		Actions: []string{statement.GetStatementCode(), statement.GetStatementCode()},
	}

	report.addToReport(statement)

	assert.Equal(t, wantReport, report)
}

func TestReportHeader(t *testing.T) {
	wantReportHeader := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n"

	gotReportHeader := ReportHeader(false)

	assert.Equal(t, wantReportHeader, gotReportHeader)
}

func TestReportHeader_WhenSafeMode(t *testing.T) {
	wantReportHeader := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report** (Reviewpad ran in dry-run mode because configuration has changed)\n\n"

	gotReportHeader := ReportHeader(true)

	assert.Equal(t, wantReportHeader, gotReportHeader)
}

func TestBuildReport(t *testing.T) {
	report := Report{
		Actions: []string{"$addLabel(\"test\")"},
	}

	wantReport := `<!--@annotation-reviewpad-report-->
**Reviewpad Report**

:scroll: **Executed actions**
` + "```yaml\n$addLabel(\"test\")\n```\n"

	gotReport := buildReport(false, &report)

	assert.Equal(t, wantReport, gotReport)
}

func TestBuildVerboseReport_WhenNoReportProvided(t *testing.T) {
	var emptyReport *Report

	wantReport := ""

	gotReport := BuildVerboseReport(emptyReport)

	assert.Equal(t, wantReport, gotReport)
}

func TestBuildVerboseReport_WhenIsProvidedReportWithNoWorkflowDetails(t *testing.T) {
	reportWithNoWorkflowDetails := &Report{}

	wantReport := ":scroll: **Executed actions**\n```yaml\n```\n"

	gotReport := BuildVerboseReport(reportWithNoWorkflowDetails)

	assert.Equal(t, wantReport, gotReport)
}

func TestBuildVerboseReport(t *testing.T) {
	report := Report{
		Actions: []string{"$addLabel(\"test\")"},
	}

	wantReport := ":scroll: **Executed actions**\n```yaml\n$addLabel(\"test\")\n```\n"

	gotReport := BuildVerboseReport(&report)

	assert.Equal(t, wantReport, gotReport)
}

func TestDeleteReportComment_WhenCommentCannotBeDeleted(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "DeleteCommentRequestFailed"
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	testCommentId := int64(1234)

	err := DeleteReportComment(mockedEnv, testCommentId)

	assert.EqualError(t, err, fmt.Sprintf("[report] error on deleting report comment %v", failMessage))
}

func TestDeleteReportComment_WhenCommentCanBeDeleted(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var deletedComment int64
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					deletedCommentId, err := strconv.Atoi(vars["comment_id"])
					if err != nil {
						assert.FailNow(t, "Delete comment request returned unexpected error: %v", err)
					}
					deletedComment = int64(deletedCommentId)
				}),
			),
		},
		nil,
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	commentToBeDeleted := int64(1234)

	err := DeleteReportComment(mockedEnv, commentToBeDeleted)

	assert.Nil(t, err)
	assert.Equal(t, commentToBeDeleted, deletedComment)
}

func TestUpdateReportComment_WhenCommentCannotBeEdited(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "EditCommentRequestFailed"
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	testCommentId := int64(1234)
	wantUpdatedComment := "Test update report comment"

	err := UpdateReportComment(mockedEnv, testCommentId, wantUpdatedComment)

	assert.EqualError(t, err, fmt.Sprintf("[report] error on updating report comment %v", failMessage))
}

func TestUpdateReportComment_WhenCommentCanBeEdited(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var gotUpdatedComment string
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := EditCommentRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotUpdatedComment = body.Body
				}),
			),
		},
		nil,
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	testCommentId := int64(1234)
	wantUpdatedComment := "Test update report comment"

	err := UpdateReportComment(mockedEnv, testCommentId, wantUpdatedComment)

	assert.Nil(t, err)
	assert.Equal(t, wantUpdatedComment, gotUpdatedComment)
}

func TestAddReportComment_WhenCommentCannotBeCreated(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "CreateCommentRequestFailed"
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	comment := "Test add report comment"

	err := AddReportComment(mockedEnv, comment)

	assert.EqualError(t, err, fmt.Sprintf("[report] error on creating report comment %v", failMessage))
}

func TestAddReportComment_WhenCommentCanBeCreated(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var createdComment string
	commentToBeCreated := "Test add report comment"

	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := EditCommentRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					createdComment = body.Body
				}),
			),
		},
		nil,
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	err := AddReportComment(mockedEnv, commentToBeCreated)

	assert.Nil(t, err)
	assert.Equal(t, createdComment, commentToBeCreated)
}

func TestFindReportComment_WhenPullRequestCommentsListingFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListCommentsRequestFailed"
	mockedEnv := MockDefaultEnv(
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	gotComment, err := FindReportComment(mockedEnv)

	assert.Nil(t, gotComment)
	assert.EqualError(t, err, fmt.Sprintf("[report] error getting issues %v", failMessage))
}

func TestFindReportComment_WhenThereIsReviewpadComment(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantComment := &github.IssueComment{
		Body: github.String("<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"),
	}
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					wantComment,
				},
			),
		},
		nil,
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	gotComment, err := FindReportComment(mockedEnv)

	assert.Nil(t, err)
	assert.Equal(t, wantComment, gotComment)
}

func TestFindReportComment_WhenThereIsNoReviewpadComment(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	comment := &github.IssueComment{
		Body: github.String("Test comment"),
	}
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					comment,
				},
			),
		},
		nil,
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	gotComment, err := FindReportComment(mockedEnv)

	assert.Nil(t, err)
	assert.Nil(t, gotComment)
}
