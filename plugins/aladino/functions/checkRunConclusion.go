// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"log"

	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func CheckRunConclusion() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           checkRunConclusionCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
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

	lastCommitSha := ghCommits[len(ghCommits)-1].GetSHA()

	checkRuns, err := e.GetGithubClient().GetCheckRunsForRef(e.GetCtx(), owner, repo, number, lastCommitSha, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, err
	}

	for _, check := range checkRuns {
		log.Printf("check: %+v", check)
		log.Printf("check-name: %+v", *check.Name)
		log.Printf("checkNname: %+v", checkRunName)
		if *check.Name == checkRunName {
			return aladino.BuildStringValue(*check.Conclusion), nil
		}
	}

	return aladino.BuildStringValue(""), nil
}
