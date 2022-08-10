// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func HasLinearHistory() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasLinearHistoryCode,
	}
}

func hasLinearHistoryCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	ghCommits, err := e.GetGithubClient().GetPullRequestCommits(e.GetCtx(), owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	for _, ghCommit := range ghCommits {
		if len(ghCommit.Parents) > 1 {
			return aladino.BuildBoolValue(false), nil
		}
	}

	return aladino.BuildBoolValue(true), nil
}
