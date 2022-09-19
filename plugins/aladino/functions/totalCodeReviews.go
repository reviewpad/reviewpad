// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/shurcooL/githubv4"
)

func TotalCodeReviews() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildIntType()),
		Code:           totalCodeReviewsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func totalCodeReviewsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	username := args[0].(*aladino.StringValue).Val

	var totalPullRequestReviewsQuery struct {
		User struct {
			ContributionsCollection struct {
				TotalPullRequestReviewContributions githubv4.Int
			} `graphql:"contributionsCollection"`
		} `graphql:"user(login: $username)"`
	}

	varGQLTotalPullRequestReviewsQuery := map[string]interface{}{
		"username": githubv4.String(username),
	}

	err := e.GetGithubClient().GetClientGraphQL().Query(e.GetCtx(), &totalPullRequestReviewsQuery, varGQLTotalPullRequestReviewsQuery)
	if err != nil {
		return nil, err
	}

	return aladino.BuildIntValue(int(totalPullRequestReviewsQuery.User.ContributionsCollection.TotalPullRequestReviewContributions)), nil
}
