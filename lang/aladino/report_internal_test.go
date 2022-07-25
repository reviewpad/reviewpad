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

	"github.com/google/go-github/v42/github"
	"github.com/gorilla/mux"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
)

type EditCommentRequestPostBody struct {
	Body string `json:"body"`
}

func TestReport_ExpectErrorWithTestingErrorMsg(t *testing.T) {
	errorMsg := "Testing Error"

	wantErr := fmt.Sprintf("[report] %v", errorMsg)
	gotErr := reportError(errorMsg)

	assert.EqualError(t, gotErr, wantErr)
}

func TestMergeReportWorkflowDetails_ExpectReportWorkflowWithTwoRules(t *testing.T) {
	left := ReportWorkflowDetails{
		Name:        "test-workflow",
		Description: "Test workflow",
		Rules: map[string]bool{
			"test-rule": true,
		},
		Actions: []string{
			"$addLabel(\"test\")",
		},
	}
	right := ReportWorkflowDetails{
		Name:        "test-workflow",
		Description: "Test workflow",
		Rules: map[string]bool{
			"tautology": true,
		},
		Actions: []string{
			"$addLabel(\"test\")",
		},
	}

	gotReportWorkflow := mergeReportWorkflowDetails(left, right)

	wantReportWorkflow := ReportWorkflowDetails{
		Name:        "test-workflow",
		Description: "Test workflow",
		Rules: map[string]bool{
			"test-rule": true,
			"tautology": true,
		},
		Actions: []string{
			"$addLabel(\"test\")",
		},
	}

	assert.Equal(t, wantReportWorkflow, gotReportWorkflow)
}

func TestAddToReport_WhenWorkflowIsNonExisting_ExpectNewWorkflowInReport(t *testing.T) {
	statement := engine.Statement{
		Code: "$addLabel(\"test\")",
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name:        "new-test-workflow",
				Description: "Testing workflow",
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: "test-rule"},
			},
		},
	}

	testWorkflow := ReportWorkflowDetails{
		Name:        statement.Metadata.Workflow.Name,
		Description: statement.Metadata.Workflow.Description,
		Rules:       map[string]bool{"tautology": true},
		Actions:     []string{statement.Code},
	}
	report := Report{
		WorkflowDetails: map[string]ReportWorkflowDetails{
			"test-workflow": testWorkflow,
		},
	}

	wantReport := Report{
		WorkflowDetails: map[string]ReportWorkflowDetails{
			"test-workflow": testWorkflow,
			"new-test-workflow": {
				Name:        statement.Metadata.Workflow.Name,
				Description: statement.Metadata.Workflow.Description,
				Rules:       map[string]bool{"test-rule": true},
				Actions:     []string{statement.Code},
			},
		},
	}

	report.addToReport(&statement)

	assert.Equal(t, wantReport, report)
}

func TestAddToReport_WhenWorkflowAlreadyExists_ExpectWorkflowWithAdditionalRule(t *testing.T) {
	statement := engine.Statement{
		Code: "$addLabel(\"test\")",
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name:        "test-workflow",
				Description: "Testing workflow",
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: "test-rule"},
			},
		},
	}
	report := Report{
		WorkflowDetails: map[string]ReportWorkflowDetails{
			"test-workflow": {
				Name:        statement.Metadata.Workflow.Name,
				Description: statement.Metadata.Workflow.Description,
				Rules:       map[string]bool{"tautology": true},
				Actions:     []string{statement.Code},
			},
		},
	}

	wantReport := Report{
		WorkflowDetails: map[string]ReportWorkflowDetails{
			"test-workflow": {
				Name:        statement.Metadata.Workflow.Name,
				Description: statement.Metadata.Workflow.Description,
				Rules: map[string]bool{
					"tautology": true,
					"test-rule": true,
				},
				Actions: []string{statement.Code},
			},
		},
	}

	report.addToReport(&statement)

	assert.Equal(t, wantReport, report)
}

func TestReportHeader_ExpectReviewpadReportHeader(t *testing.T) {
	wantReportHeader := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n"

	gotReportHeader := ReportHeader()

	assert.Equal(t, wantReportHeader, gotReportHeader)
}

func TestBuildReport_ExpectReportWithTestWorkflowDetails(t *testing.T) {
	report := Report{
		WorkflowDetails: map[string]ReportWorkflowDetails{
			"test-workflow": {
				Name:        "test-workflow",
				Description: "Testing workflow",
				Rules:       map[string]bool{"tautology": true},
				Actions:     []string{"$addLabel(\"test\")"},
			},
		},
	}

	wantReport := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\n| Workflows <sub><sup>activated</sup></sub> | Rules <sub><sup>triggered</sup></sub> | Actions <sub><sup>ran</sub></sup> | Description |\n| - | - | - | - |\n| test-workflow | tautology<br> | `$addLabel(\"test\")`<br> | Testing workflow |\n"

	gotReport := buildReport(&report)

	assert.Equal(t, wantReport, gotReport)
}

func TestBuildVerboseReport_WhenNoReportProvided_ExpectEmptyReport(t *testing.T) {
	var emptyReport *Report

	wantReport := ""

	gotReport := BuildVerboseReport(emptyReport)

	assert.Equal(t, wantReport, gotReport)
}

func TestBuildVerboseReport_WhenIsProvidedReportWithNoWorkflowDetails_ExpectReportWithNoWorkflowsActivated(t *testing.T) {
	reportWithNoWorkflowDetails := &Report{}

	wantReport := ":scroll: **Explanation**\nNo workflows activated"

	gotReport := BuildVerboseReport(reportWithNoWorkflowDetails)

	assert.Equal(t, wantReport, gotReport)
}

func TestBuildVerboseReport_ExpectReportWithTestWorkflowDetails(t *testing.T) {
	report := Report{
		WorkflowDetails: map[string]ReportWorkflowDetails{
			"test-workflow": {
				Name:        "test-workflow",
				Description: "Testing workflow",
				Rules:       map[string]bool{"tautology": true},
				Actions:     []string{"$addLabel(\"test\")"},
			},
		},
	}

	wantReport := ":scroll: **Explanation**\n| Workflows <sub><sup>activated</sup></sub> | Rules <sub><sup>triggered</sup></sub> | Actions <sub><sup>ran</sub></sup> | Description |\n| - | - | - | - |\n| test-workflow | tautology<br> | `$addLabel(\"test\")`<br> | Testing workflow |\n"

	gotReport := BuildVerboseReport(&report)

	assert.Equal(t, wantReport, gotReport)
}

func TestDeleteReportComment_ExpectDeleteCommetRequestFail(t *testing.T) {
	failMessage := "DeleteCommentRequestFailed"
	mockedEnv, err := MockDefaultEnv(
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
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	testCommentId := int64(1234)

	err = DeleteReportComment(mockedEnv, testCommentId)

	assert.EqualError(t, err, fmt.Sprintf("[report] error on deleting report comment %v", failMessage))
}

func TestDeleteReportComment_ExpectNoError(t *testing.T) {
	var deletedComment int64
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					deletedCommentId, err := strconv.Atoi(vars["comment_id"])
					if err != nil {
						assert.Fail(t, "Delete comment request returned unexpected error: %v", err)
					}
					deletedComment = int64(deletedCommentId)
				}),
			),
		},
		nil,
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	commentToBeDeleted := int64(1234)

	err = DeleteReportComment(mockedEnv, commentToBeDeleted)

	assert.Nil(t, err)
	assert.Equal(t, commentToBeDeleted, deletedComment)
}

func TestUpdateReportComment_ExpectEditCommentRequestFail(t *testing.T) {
	failMessage := "EditCommentRequestFailed"
	mockedEnv, err := MockDefaultEnv(
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
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	testCommentId := int64(1234)
	wantUpdatedComment := "Test update report comment"

	err = UpdateReportComment(mockedEnv, testCommentId, wantUpdatedComment)

	assert.EqualError(t, err, fmt.Sprintf("[report] error on updating report comment %v", failMessage))
}

func TestUpdateReportComment_ExpectNoError(t *testing.T) {
	var gotUpdatedComment string
	mockedEnv, err := MockDefaultEnv(
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
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	testCommentId := int64(1234)
	wantUpdatedComment := "Test update report comment"

	err = UpdateReportComment(mockedEnv, testCommentId, wantUpdatedComment)

	assert.Nil(t, err)
	assert.Equal(t, wantUpdatedComment, gotUpdatedComment)
}

func TestAddReportComment_ExpectCreateCommentRequestFail(t *testing.T) {
	failMessage := "CreateCommentRequestFailed"
	mockedEnv, err := MockDefaultEnv(
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
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	comment := "Test add report comment"

	err = AddReportComment(mockedEnv, comment)

	assert.EqualError(t, err, fmt.Sprintf("[report] error on creating report comment %v", failMessage))
}

func TestAddReportComment_ExpectNoError(t *testing.T) {
	var createdComment string
	commentToBeCreated := "Test add report comment"

	mockedEnv, err := MockDefaultEnv(
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
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	err = AddReportComment(mockedEnv, commentToBeCreated)

	assert.Nil(t, err)
	assert.Equal(t, createdComment, commentToBeCreated)
}

func TestFindReportComment_ExpectGetPullRequestCommentsFail(t *testing.T) {
	failMessage := "ListCommentsRequestFailed"
	mockedEnv, err := MockDefaultEnv(
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
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	gotComment, err := FindReportComment(mockedEnv)

	assert.Nil(t, gotComment)
	assert.EqualError(t, err, fmt.Sprintf("[report] error getting issues %v", failMessage))
}

func TestFindReportComment_ExpectReviewpadComment(t *testing.T) {
	wantComment := &github.IssueComment{
		Body: github.String("<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"),
	}
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					wantComment,
				},
			),
		},
		nil,
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	gotComment, err := FindReportComment(mockedEnv)

	assert.Nil(t, err)
	assert.Equal(t, wantComment, gotComment)
}

func TestFindReportComment_ExpectNoReviewpadComment(t *testing.T) {
	comment := &github.IssueComment{
		Body: github.String("Test comment"),
	}
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					comment,
				},
			),
		},
		nil,
	)
	if err != nil {
		assert.Fail(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	gotComment, err := FindReportComment(mockedEnv)

	assert.Nil(t, err)
	assert.Nil(t, gotComment)
}
