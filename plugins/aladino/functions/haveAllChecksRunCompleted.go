// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HaveAllChecksRunCompleted() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: lang.BuildFunctionType([]lang.Type{
			lang.BuildArrayOfType(lang.BuildStringType()),
			lang.BuildStringType(),
			lang.BuildArrayOfType(lang.BuildStringType()),
		}, lang.BuildBoolType()),
		Code:           haveAllChecksRunCompleted,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func haveAllChecksRunCompleted(e aladino.Env, args []lang.Value) (lang.Value, error) {
	checkRunsToIgnore := args[0].(*lang.ArrayValue)
	conclusion := args[1].(*lang.StringValue)
	checkConclusionsToIgnore := args[2].(*lang.ArrayValue)
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	owner := pullRequest.GetTargetEntity().Owner
	repo := pullRequest.GetTargetEntity().Repo
	number := pullRequest.GetTargetEntity().Number

	lastCommitSHA, err := pullRequest.GetLastCommit()
	if err != nil {
		return nil, err
	}

	if lastCommitSHA == "" {
		return lang.BuildBoolValue(false), nil
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
		ignoredRuns[item.(*lang.StringValue).Val] = true
	}

	ignoredConclusions := map[string]bool{}
	for _, ignoredConclusion := range checkConclusionsToIgnore.Vals {
		ignoredConclusions[ignoredConclusion.(*lang.StringValue).Val] = true
	}

	for _, checkRun := range checkRuns {
		if isIgnoredCheckRun(checkRun, ignoredRuns, ignoredConclusions) {
			continue
		}

		if checkRun.GetStatus() != "completed" {
			return lang.BuildBoolValue(false), nil
		}

		if conclusion.Val != "" && checkRun.GetConclusion() != conclusion.Val {
			return lang.BuildBoolValue(false), nil
		}
	}

	return lang.BuildBoolValue(true), nil
}

func isIgnoredCheckRun(checkRun *github.CheckRun, ignoredRuns, ignoredConclusions map[string]bool) bool {
	return ignoredRuns[checkRun.GetName()] || ignoredConclusions[checkRun.GetConclusion()]
}
