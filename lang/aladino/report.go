// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
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

const ReviewpadReportCommentAnnotation = "<!--@annotation-reviewpad-report-->"

func reportError(format string, a ...interface{}) error {
	return fmtio.Errorf("report", format, a...)
}

func mergeReportWorkflowDetails(left, right ReportWorkflowDetails) ReportWorkflowDetails {
	for rule := range right.Rules {
		if _, ok := left.Rules[rule]; !ok {
			left.Rules[rule] = true
		}
	}

	return left
}

func (r *Report) SendReport(ctx context.Context, reviewpadFileChanged bool, mode string, pr *github.PullRequest, client *github.Client) error {
	execLog("generating report")

	var err error

	comment, err := FindReportComment(ctx, pr, client)
	if err != nil {
		return err
	}

	if mode == engine.SILENT_MODE {
		if comment != nil {
			return DeleteReportComment(ctx, pr, *comment.ID, client)
		}
		return nil
	}

	report := buildReport(reviewpadFileChanged, r)

	if comment == nil {
		return AddReportComment(ctx, pr, report, client)
	}

	return UpdateReportComment(ctx, pr, *comment.ID, report, client)
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

func ReportFromProgram(program *engine.Program) *Report {
	report := &Report{
		WorkflowDetails: make(map[string]ReportWorkflowDetails),
	}

	for _, statement := range program.Statements {
		report.addToReport(statement)
	}

	return report
}

func ReportHeader(reviewpadFileChanged bool) string {
	var sb strings.Builder

	// Annotation
	sb.WriteString(fmt.Sprintf("%v\n", ReviewpadReportCommentAnnotation))
	// Header
	if reviewpadFileChanged {
		sb.WriteString("**Reviewpad Report** (Reviewpad ran in dry-run mode because configuration has changed)\n\n")
	} else {
		sb.WriteString("**Reviewpad Report**\n\n")
	}

	return sb.String()
}

func buildReport(reviewpadFileChanged bool, report *Report) string {
	var sb strings.Builder

	sb.WriteString(ReportHeader(reviewpadFileChanged))
	sb.WriteString(BuildVerboseReport(report))

	return sb.String()
}

func BuildVerboseReport(report *Report) string {
	if report == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(":scroll: **Explanation**\n")

	reportDetails := report.WorkflowDetails

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
		for actRule := range workflow.Rules {
			actRules += fmt.Sprintf("%v<br>", actRule)
		}

		actActions := ""
		for _, actAction := range workflow.Actions {
			actActions += fmt.Sprintf("`%v`<br>", actAction)
		}

		sb.WriteString(fmt.Sprintf("| %v | %v | %v | %v |\n", workflow.Name, actRules, actActions, workflow.Description))
	}

	return sb.String()
}

func DeleteReportComment(ctx context.Context, pr *github.PullRequest, commentId int64, client *github.Client) error {
	owner := utils.GetPullRequestBaseOwnerName(pr)
	repo := utils.GetPullRequestBaseRepoName(pr)

	_, err := client.Issues.DeleteComment(ctx, owner, repo, commentId)

	if err != nil {
		return reportError("error on deleting report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func UpdateReportComment(ctx context.Context, pr *github.PullRequest, commentId int64, report string, client *github.Client) error {
	gitHubComment := github.IssueComment{
		Body: &report,
	}

	owner := utils.GetPullRequestBaseOwnerName(pr)
	repo := utils.GetPullRequestBaseRepoName(pr)

	_, _, err := client.Issues.EditComment(ctx, owner, repo, commentId, &gitHubComment)

	if err != nil {
		return reportError("error on updating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func AddReportComment(ctx context.Context, pr *github.PullRequest, report string, client *github.Client) error {
	gitHubComment := github.IssueComment{
		Body: &report,
	}

	owner := utils.GetPullRequestBaseOwnerName(pr)
	repo := utils.GetPullRequestBaseRepoName(pr)
	prNum := utils.GetPullRequestNumber(pr)

	_, _, err := client.Issues.CreateComment(ctx, owner, repo, prNum, &gitHubComment)

	if err != nil {
		return reportError("error on creating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func FindReportComment(ctx context.Context, pr *github.PullRequest, client *github.Client) (*github.IssueComment, error) {
	owner := utils.GetPullRequestBaseOwnerName(pr)
	repo := utils.GetPullRequestBaseRepoName(pr)
	prNum := utils.GetPullRequestNumber(pr)

	comments, err := utils.GetPullRequestComments(ctx, client, owner, repo, prNum, &github.IssueListCommentsOptions{
		Sort:      github.String("created"),
		Direction: github.String("asc"),
	})
	if err != nil {
		return nil, reportError("error getting issues %v", err.(*github.ErrorResponse).Message)
	}

	reviewpadCommentAnnotationRegex := regexp.MustCompile(fmt.Sprintf("^%v", ReviewpadReportCommentAnnotation))

	var reviewpadExistingComment *github.IssueComment

	for _, comment := range comments {
		isReviewpadReportComment := reviewpadCommentAnnotationRegex.Match([]byte(*comment.Body))
		if isReviewpadReportComment {
			reviewpadExistingComment = comment
			break
		}
	}

	return reviewpadExistingComment, nil
}
