// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/google/go-github/v49/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
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
		return nil, fmt.Errorf("failed to get last commit: %s", err.Error())
	}

	if lastCommitSHA == "" {
		return aladino.BuildBoolValue(false), nil
	}

	checkRuns, err := e.GetGithubClient().GetCheckRunsForRef(e.GetCtx(), owner, repo, number, lastCommitSHA, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get check runs: %s", err.Error())
	}

	checkRunIgnored := map[string]bool{}
	for _, item := range checkRunsToIgnore.Vals {
		checkRunIgnored[item.(*aladino.StringValue).Val] = true
	}

	checkConclusionsConsidered := map[string]bool{}
	for _, ignoredConclusion := range checkConclusions.Vals {
		checkConclusionsConsidered[ignoredConclusion.(*aladino.StringValue).Val] = true
	}

	for _, checkRun := range checkRuns {
		if checkRunIgnored[checkRun.GetName()] {
			continue
		}

		if checkRun.GetStatus() != "completed" {
			continue
		}

		// If no conclusions are specified, we consider all check runs as valid
		if len(checkConclusionsConsidered) == 0 {
			return aladino.BuildBoolValue(true), nil
		}

		if checkConclusionsConsidered[checkRun.GetConclusion()] {
			return aladino.BuildBoolValue(true), nil
		}
	}

	return aladino.BuildBoolValue(false), nil
}
