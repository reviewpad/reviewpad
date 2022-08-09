// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func HasUnaddressedThreads() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasUnaddressedThreadsCode,
	}
}

func hasUnaddressedThreadsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	totalRequestTries := 2

	reviewThreads, err := e.GetGithubClient().GetReviewThreads(e.GetCtx(), owner, repo, prNum, totalRequestTries)
	if err != nil {
		return nil, err
	}

	for _, reviewThread := range reviewThreads {
		if !reviewThread.IsResolved && !reviewThread.IsOutdated {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}
