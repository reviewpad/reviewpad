// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HaveAllChecksRunCompleted() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{
			aladino.BuildArrayOfType(aladino.BuildStringType()),
			aladino.BuildStringType(),
			aladino.BuildArrayOfType(aladino.BuildStringType()),
		}, aladino.BuildBoolType()),
		Code:           haveAllChecksRunCompleted,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func haveAllChecksRunCompleted(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	checkRunsToIgnore := args[0].(*aladino.ArrayValue)
	conclusion := args[1].(*aladino.StringValue)
	checkConclusionsToIgnore := args[2].(*aladino.ArrayValue)
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

	ignoredRuns := map[string]bool{
		// By default, ignore Reviewpad's own check runs
		"reviewpad": true,
	}
	for _, item := range checkRunsToIgnore.Vals {
		ignoredRuns[item.(*aladino.StringValue).Val] = true
	}

	ignoredConclusions := map[string]bool{}
	for _, ignoredConclusion := range checkConclusionsToIgnore.Vals {
		ignoredConclusions[ignoredConclusion.(*aladino.StringValue).Val] = true
	}

	for _, checkRun := range checkRuns {
		if isIgnoredCheckRun(checkRun, ignoredRuns, ignoredConclusions) {
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

func isIgnoredCheckRun(checkRun *github.CheckRun, ignoredRuns, ignoredConclusions map[string]bool) bool {
	return ignoredRuns[checkRun.GetName()] || ignoredConclusions[checkRun.GetConclusion()]
}
