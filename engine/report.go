// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/utils"
	"github.com/reviewpad/reviewpad/utils/fmtio"
)

type ReportWorkflowDetails struct {
	name            string
	description     string
	actRules        []string
	actActions      []string
	actExtraActions []string
}

var reviewpadReportCommentAnnotation = "<!--@annotation-reviewpad-->"

func reportError(format string, a ...interface{}) error {
	return fmtio.Errorf("report", format, a...)
}

func buildReport(reportDetails *[]ReportWorkflowDetails) string {
	var sb strings.Builder

	// Annotation
	sb.WriteString(fmt.Sprintf("%v\n", reviewpadReportCommentAnnotation))
	// Header
	sb.WriteString("### Reviewpad report\n")
	sb.WriteString("___\n")

	if len(*reportDetails) == 0 {
		sb.WriteString("No workflows activated")
		msg := sb.String()
		return msg
	}

	// Report
	sb.WriteString("| Workflows <sub><sup>activated</sup></sub> | Rules <sub><sup>triggered</sup></sub> | Actions <sub><sup>ran</sub></sup> | Description |\n")
	sb.WriteString("| - | - | - | - |\n")

	for _, workflow := range *reportDetails {
		actRules := ""
		for _, actRule := range workflow.actRules {
			actRules += fmt.Sprintf("%v<br>", actRule)
		}

		actActions := ""
		for _, actAction := range workflow.actActions {
			actActions += fmt.Sprintf("`%v`<br>", actAction)
		}
		for _, actExtraAction := range workflow.actExtraActions {
			actActions += fmt.Sprintf("[e] `%v`<br>", actExtraAction)
		}

		sb.WriteString(fmt.Sprintf("| %v | %v | %v | %v |\n", workflow.name, actRules, actActions, workflow.description))
	}

	msg := sb.String()
	return msg
}

func updateReportComment(env *Env, commentId int64, report string) error {
	gitHubComment := github.IssueComment{
		Body: &report,
	}

	owner := utils.GetPullRequestOwnerName(env.PullRequest)
	repo := utils.GetPullRequestRepoName(env.PullRequest)

	_, _, err := env.Client.Issues.EditComment(env.Ctx, owner, repo, commentId, &gitHubComment)

	if err != nil {
		return reportError("error on updating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func addReportComment(env *Env, prNum int, report string) error {
	gitHubComment := github.IssueComment{
		Body: &report,
	}

	owner := utils.GetPullRequestOwnerName(env.PullRequest)
	repo := utils.GetPullRequestRepoName(env.PullRequest)

	_, _, err := env.Client.Issues.CreateComment(env.Ctx, owner, repo, prNum, &gitHubComment)

	if err != nil {
		return reportError("error on creating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return nil
}

func reportProgram(env *Env, reportDetails *[]ReportWorkflowDetails) (string, error) {
	owner := utils.GetPullRequestOwnerName(env.PullRequest)
	repo := utils.GetPullRequestRepoName(env.PullRequest)
	prNum := utils.GetPullRequestNumber(env.PullRequest)

	comments, _, err := env.Client.Issues.ListComments(env.Ctx, owner, repo, prNum, &github.IssueListCommentsOptions{
		Sort:      github.String("created"),
		Direction: github.String("asc"),
	})

	if err != nil {
		return "", reportError("error getting issues %v", err.(*github.ErrorResponse).Message)
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
		return "", err
	}

	return report, nil
}
