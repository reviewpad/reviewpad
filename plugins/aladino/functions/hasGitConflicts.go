// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/shurcooL/githubv4"
)

func HasGitConflicts() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           hasGitConflictsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasGitConflictsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).CodeReview

	prNum := host.GetCodeReviewNumber(pullRequest)
	repoOwner := host.GetCodeReviewBaseOwnerName(pullRequest)
	repoName := host.GetCodeReviewBaseRepoName(pullRequest)
	var pullRequestQuery struct {
		Repository struct {
			PullRequest struct {
				Mergeable githubv4.String
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGQLPullRequestQuery := map[string]interface{}{
		"pullRequestNumber": githubv4.Int(prNum),
		"repositoryOwner":   githubv4.String(repoOwner),
		"repositoryName":    githubv4.String(repoName),
	}

	err := e.GetGithubClient().GetClientGraphQL().Query(e.GetCtx(), &pullRequestQuery, varGQLPullRequestQuery)
	if err != nil {
		return nil, err
	}

	if string(pullRequestQuery.Repository.PullRequest.Mergeable) == "CONFLICTING" {
		return aladino.BuildBoolValue(true), nil
	}

	return aladino.BuildBoolValue(false), nil
}
