// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func CheckRunStatus() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           checkRunStatusCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func checkRunStatusCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	checkRunName := args[0].(*aladino.StringValue).Val
	entity := e.GetTarget().GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo
	number := entity.Number

	commits, _, err := e.GetGithubClient().GetClientREST().PullRequests.ListCommits(e.GetCtx(), owner, repo, number, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	lastCommitSha := commits[len(commits)-1].GetSHA()

	checkRuns, _, err := e.GetGithubClient().ListCheckRunsForRef(e.GetCtx(), owner, repo, lastCommitSha, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, err
	}

	for _, check := range checkRuns.CheckRuns {
		if *check.Name == checkRunName {
			return aladino.BuildStringValue(*check.Status), nil
		}
	}

	return aladino.BuildStringValue(""), nil
}
