// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/shurcooL/githubv4"
)

func HasLinkedIssues() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasLinkedIssuesCode,
	}
}

func hasLinkedIssuesCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	var pullRequestQuery struct {
		Repository struct {
			PullRequest struct {
				ClosingIssuesReferences struct {
					TotalCount githubv4.Int
				}
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	pullRequest := e.GetPullRequest()
	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	varGQLPullRequestQuery := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(prNum),
	}

	// FIXME: Move to codehost/github package
	err := e.GetGithubClient().GetClientGraphQL().Query(e.GetCtx(), &pullRequestQuery, varGQLPullRequestQuery)

	if err != nil {
		return nil, err
	}

	return aladino.BuildBoolValue(pullRequestQuery.Repository.PullRequest.ClosingIssuesReferences.TotalCount > 0), nil
}
