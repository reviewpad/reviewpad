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

func CheckRunConclusion() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           checkRunConclusionCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func checkRunConclusionCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	checkRunName := args[0].(*aladino.StringValue).Val
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	owner := pullRequest.GetTargetEntity().Owner
	repo := pullRequest.GetTargetEntity().Repo
	number := pullRequest.GetTargetEntity().Number

	ghCommits, err := e.GetGithubClient().GetPullRequestCommits(e.GetCtx(), owner, repo, number)
	if err != nil {
		return nil, err
	}

	if len(ghCommits) == 0 {
		return aladino.BuildStringValue(""), nil
	}

	lastCommitSha := ghCommits[len(ghCommits)-1].GetSHA()

	checkRuns, err := e.GetGithubClient().GetCheckRunsForRef(e.GetCtx(), owner, repo, number, lastCommitSha, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, err
	}

	var checkRunConclusion string
	var checkRunCompletedAt *github.Timestamp

	for _, checkRun := range checkRuns {
		if isCheckEligible(checkRun, checkRunName, checkRunCompletedAt) {
			checkRunConclusion = *checkRun.Conclusion
			checkRunCompletedAt = checkRun.CompletedAt
		}
	}

	return aladino.BuildStringValue(checkRunConclusion), nil
}

func isCheckEligible(checkRun *github.CheckRun, requiredName string, requiredMinCompletedAt *github.Timestamp) bool {
	if *checkRun.Name != requiredName {
		return false
	}

	if *checkRun.Status != "completed" {
		return false
	}

	if requiredMinCompletedAt != nil && checkRun.CompletedAt.Before(requiredMinCompletedAt.Time) {
		return false
	}

	if checkRun.Conclusion == nil {
		return false
	}

	return true
}
