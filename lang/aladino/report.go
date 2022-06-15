// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/engine"
	"github.com/reviewpad/reviewpad/utils"
	"github.com/reviewpad/reviewpad/utils/fmtio"
)

type Report struct {
	WorkflowDetails map[string]ReportWorkflowDetails
}

type ReportWorkflowDetails struct {
	Name        string
	Description string
	Rules       map[string]bool
	Actions     []string
}

var reviewpadReportCommentAnnotation = "<!--@annotation-reviewpad-->"

func reportError(format string, a ...interface{}) error {
	return fmtio.Errorf("report", format, a...)
}

func mergeReportWorkflowDetails(left, right ReportWorkflowDetails) ReportWorkflowDetails {
	for rule, _ := range right.Rules {
		if _, ok := left.Rules[rule]; !ok {
			left.Rules[rule] = true
		}
	}

	return left
}

func (report *Report) addToReport(statement *engine.Statement) {
	workflowName := statement.Metadata.Workflow.Name

	rules := make(map[string]bool, len(statement.Metadata.TriggeredBy))
	for _, rule := range statement.Metadata.TriggeredBy {
		rules[rule.Rule] = true
	}

	reportWorkflow := ReportWorkflowDetails{
		Name:        workflowName,
		Description: statement.Metadata.Workflow.Description,
		Rules:       rules,
		Actions:     []string{statement.Code},
	}

	workflow, ok := report.WorkflowDetails[workflowName]
	if !ok {
		report.WorkflowDetails[workflowName] = reportWorkflow
	} else {
		report.WorkflowDetails[workflowName] = mergeReportWorkflowDetails(workflow, reportWorkflow)
	}
}

func buildReport(reportDetails map[string]ReportWorkflowDetails) string {
	var sb strings.Builder

	// Annotation
	sb.WriteString(fmt.Sprintf("%v\n", reviewpadReportCommentAnnotation))
	// Header
	sb.WriteString("### Reviewpad report\n")
	sb.WriteString("___\n")

	if len(reportDetails) == 0 {
		sb.WriteString("No workflows activated")
		msg := sb.String()
		return msg
	}

	// Report
	sb.WriteString("| Workflows <sub><sup>activated</sup></sub> | Rules <sub><sup>triggered</sup></sub> | Actions <sub><sup>ran</sub></sup> | Description |\n")
	sb.WriteString("| - | - | - | - |\n")

	for _, workflow := range reportDetails {
		actRules := ""
		for _, actRule := range workflow.Rules {
			actRules += fmt.Sprintf("%v<br>", actRule)
		}

		actActions := ""
		for _, actAction := range workflow.Actions {
			actActions += fmt.Sprintf("`%v`<br>", actAction)
		}

		sb.WriteString(fmt.Sprintf("| %v | %v | %v | %v |\n", workflow.Name, actRules, actActions, workflow.Description))
	}

	msg := sb.String()
	return msg
}

func updateReportComment(env Env, commentId int64, report string) error {
	gitHubComment := github.IssueComment{
		Body: &report,
	}

	pullRequest := env.GetPullRequest()
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	_, _, err := env.GetClient().Issues.EditComment(env.GetCtx(), owner, repo, commentId, &gitHubComment)

	if err != nil {
		return reportError("error on updating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func addReportComment(env Env, prNum int, report string) error {
	gitHubComment := github.IssueComment{
		Body: &report,
	}

	pullRequest := env.GetPullRequest()
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	_, _, err := env.GetClient().Issues.CreateComment(env.GetCtx(), owner, repo, prNum, &gitHubComment)

	if err != nil {
		return reportError("error on creating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func ReportProgram(env Env) error {
	reportDetails := env.GetReport().WorkflowDetails
	pullRequest := env.GetPullRequest()
	owner := utils.GetPullRequestOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)
	prNum := utils.GetPullRequestNumber(pullRequest)

	comments, _, err := env.GetClient().Issues.ListComments(env.GetCtx(), owner, repo, prNum, &github.IssueListCommentsOptions{
		Sort:      github.String("created"),
		Direction: github.String("asc"),
	})

	if err != nil {
		return reportError("error getting issues %v", err.(*github.ErrorResponse).Message)
	}

	reviewpadCommentAnnotationRegex := regexp.MustCompile(fmt.Sprintf("^%v", reviewpadReportCommentAnnotation))

	var reviewpadExistingComment *github.IssueComment

	for _, comment := range comments {
		isReviewpadComment := reviewpadCommentAnnotationRegex.Match([]byte(*comment.Body))
		if isReviewpadComment {
			reviewpadExistingComment = comment
			break
		}
	}

	report := buildReport(reportDetails)

	if reviewpadExistingComment == nil {
		err = addReportComment(env, prNum, report)
	} else {
		err = updateReportComment(env, *reviewpadExistingComment.ID, report)
	}

	if err != nil {
		return err
	}

	return nil
}
