// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func FailCheckStatus() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           failCheckStatusCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func failCheckStatusCode(e aladino.Env, args []aladino.Value) error {
	targetEntity := e.GetTarget().GetTargetEntity()
	failMessage := args[0].(*aladino.StringValue).Val

	e.GetBuiltInsReportedMessages()[aladino.SEVERITY_FAIL] = append(e.GetBuiltInsReportedMessages()[aladino.SEVERITY_FAIL], failMessage)

	if e.GetCheckRunID() != nil {
		e.SetCheckRunConclusion("failure")

		summary := strings.Builder{}

		for _, msg := range e.GetBuiltInsReportedMessages()[aladino.SEVERITY_FAIL] {
			summary.WriteString(fmt.Sprintf("- %s\n", msg))
		}

		_, _, err := e.GetGithubClient().GetClientREST().Checks.UpdateCheckRun(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, *e.GetCheckRunID(), github.UpdateCheckRunOptions{
			Name:       "reviewpad",
			Status:     github.String("completed"),
			Conclusion: github.String("failure"),
			Output: &github.CheckRunOutput{
				Title:   github.String("Reviewpad policy failed"),
				Summary: github.String(summary.String()),
			},
		})
		return err
	}

	return nil
}
