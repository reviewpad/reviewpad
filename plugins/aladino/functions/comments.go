// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v45/github"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Comments() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: commentsCode,
	}
}

func commentsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	ghComments, err := e.GetGithubClient().GetPullRequestComments(e.GetCtx(), owner, repo, prNum, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	commentsBody := make([]aladino.Value, len(ghComments))
	for i, ghComment := range ghComments {
		commentsBody[i] = aladino.BuildStringValue(ghComment.GetBody())
	}

	return aladino.BuildArrayValue(commentsBody), nil
}
