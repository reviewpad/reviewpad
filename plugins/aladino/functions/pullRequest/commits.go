// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_pullRequest

import (
	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func Commits() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: commitsCode,
	}
}

func commitsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	ghCommits, _, err := e.GetClient().PullRequests.ListCommits(e.GetCtx(), owner, repo, prNum, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	commitMessages := make([]aladino.Value, len(ghCommits))
	for i, ghCommit := range ghCommits {
		commitMessages[i] = aladino.BuildStringValue(ghCommit.Commit.GetMessage())
	}

	return aladino.BuildArrayValue(commitMessages), nil
}
