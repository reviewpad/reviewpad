// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func HaveAllChecksRunCompleted() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           haveAllChecksRunCompleted,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func haveAllChecksRunCompleted(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	checkRunsToIgnore := args[0].(*aladino.ArrayValue)
	conclusion := args[1].(*aladino.StringValue)
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	owner := pullRequest.GetTargetEntity().Owner
	repo := pullRequest.GetTargetEntity().Repo
	number := pullRequest.GetTargetEntity().Number

	ghCommits, err := e.GetGithubClient().GetPullRequestCommits(e.GetCtx(), owner, repo, number)
	if err != nil {
		return nil, err
	}

	if len(ghCommits) == 0 {
		return aladino.BuildBoolValue(false), nil
	}

	lastCommitSha := ghCommits[len(ghCommits)-1].GetSHA()

	checkRuns, err := e.GetGithubClient().GetCheckRunsForRef(e.GetCtx(), owner, repo, number, lastCommitSha, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, err
	}

	ignoredRuns := map[string]bool{}
	for _, item := range checkRunsToIgnore.Vals {
		ignoredRuns[item.(*aladino.StringValue).Val] = true
	}

	for _, checkRun := range checkRuns {
		if ignoredRuns[checkRun.GetName()] {
			continue
		}

		if checkRun.GetStatus() != "completed" {
			return aladino.BuildBoolValue(false), nil
		}

		if conclusion.Val != "" && checkRun.GetConclusion() != conclusion.Val {
			return aladino.BuildBoolValue(false), nil
		}
	}

	return aladino.BuildBoolValue(true), nil
}
