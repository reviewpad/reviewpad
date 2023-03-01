// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func HasAnyCheckRunCompleted() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{
			aladino.BuildArrayOfType(aladino.BuildStringType()),
			aladino.BuildArrayOfType(aladino.BuildStringType()),
		}, aladino.BuildBoolType()),
		Code:           hasAnyCheckRunCompleted,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasAnyCheckRunCompleted(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	checkRunsToIgnore := args[0].(*aladino.ArrayValue)
	checkConclusions := args[1].(*aladino.ArrayValue)
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	owner := pullRequest.GetTargetEntity().Owner
	repo := pullRequest.GetTargetEntity().Repo
	number := pullRequest.GetTargetEntity().Number

	lastCommitSHA, err := pullRequest.GetLastCommit()
	if err != nil {
		return nil, err
	}

	if lastCommitSHA == "" {
		return aladino.BuildBoolValue(false), nil
	}

	checkRuns, err := e.GetGithubClient().GetCheckRunsForRef(e.GetCtx(), owner, repo, number, lastCommitSHA, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, err
	}

	ignoredRuns := map[string]bool{}
	for _, item := range checkRunsToIgnore.Vals {
		ignoredRuns[item.(*aladino.StringValue).Val] = true
	}

	checkConclusionsToCheck := map[string]bool{}
	for _, ignoredConclusion := range checkConclusions.Vals {
		checkConclusionsToCheck[ignoredConclusion.(*aladino.StringValue).Val] = true
	}

	for _, checkRun := range checkRuns {
		if isExcludedCheckRun(checkRun, ignoredRuns) {
			continue
		}

		if checkConclusionsToCheck[checkRun.GetConclusion()] {
			return aladino.BuildBoolValue(true), nil
		}
	}

	return aladino.BuildBoolValue(false), nil
}

func isExcludedCheckRun(checkRun *github.CheckRun, ignoredRuns map[string]bool) bool {
	return ignoredRuns[checkRun.GetName()] || checkRun.GetStatus() != "completed"
}
