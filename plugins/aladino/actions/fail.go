// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Fail() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           failCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func failCode(e aladino.Env, args []aladino.Value) error {
	targetEntity := e.GetTarget().GetTargetEntity()
	failMessage := args[0].(*aladino.StringValue).Val

	e.GetBuiltInsReportedMessages()[aladino.SEVERITY_FATAL] = append(e.GetBuiltInsReportedMessages()[aladino.SEVERITY_FATAL], failMessage)

	if e.GetCheckRunID() != nil {
		_, _, err := e.GetGithubClient().GetClientREST().Checks.UpdateCheckRun(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, *e.GetCheckRunID(), github.UpdateCheckRunOptions{
			Status:     github.String("completed"),
			Conclusion: github.String("failure"),
		})
		return err
	}

	return nil
}
