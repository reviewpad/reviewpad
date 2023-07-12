// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasAnyCheckRunCompleted() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: lang.BuildFunctionType([]lang.Type{
			lang.BuildArrayOfType(lang.BuildStringType()),
			lang.BuildArrayOfType(lang.BuildStringType()),
		}, lang.BuildBoolType()),
		Code:           hasAnyCheckRunCompleted,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func hasAnyCheckRunCompleted(e aladino.Env, args []lang.Value) (lang.Value, error) {
	checkRunsToIgnore := args[0].(*lang.ArrayValue)
	checkConclusions := args[1].(*lang.ArrayValue)

	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	owner := pullRequest.GetTargetEntity().Owner
	repo := pullRequest.GetTargetEntity().Repo
	number := pullRequest.GetTargetEntity().Number

	lastCommitSHA, err := pullRequest.GetLastCommit()
	if err != nil {
		return nil, fmt.Errorf("failed to get last commit: %s", err.Error())
	}

	if lastCommitSHA == "" {
		return lang.BuildBoolValue(false), nil
	}

	checkRuns, err := e.GetGithubClient().GetCheckRunsForRef(e.GetCtx(), owner, repo, number, lastCommitSHA, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get check runs: %s", err.Error())
	}

	checkRunIgnored := map[string]bool{
		// By default, ignore Reviewpad's own check runs
		"reviewpad": true,
	}
	for _, item := range checkRunsToIgnore.Vals {
		checkRunIgnored[item.(*lang.StringValue).Val] = true
	}

	checkConclusionsConsidered := map[string]bool{}
	for _, ignoredConclusion := range checkConclusions.Vals {
		checkConclusionsConsidered[ignoredConclusion.(*lang.StringValue).Val] = true
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
			return lang.BuildBoolValue(true), nil
		}

		if checkConclusionsConsidered[checkRun.GetConclusion()] {
			return lang.BuildBoolValue(true), nil
		}
	}

	return lang.BuildBoolValue(false), nil
}
