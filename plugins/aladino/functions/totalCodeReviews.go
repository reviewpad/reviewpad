// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/shurcooL/githubv4"
)

func TotalCodeReviews() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildIntType()),
		Code:           totalCodeReviewsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func totalCodeReviewsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	username := args[0].(*aladino.StringValue).Val

	var totalPullRequestReviewContributionsQuery struct {
		User struct {
			ContributionsCollection struct {
				TotalPullRequestReviewContributions githubv4.Int
			} `graphql:"contributionsCollection"`
		} `graphql:"user(login: $username)"`
	}

	varGQLTotalPullRequestReviewContributionsQuery := map[string]interface{}{
		"username": githubv4.String(username),
	}

	err := e.GetGithubClient().GetClientGraphQL().Query(e.GetCtx(), &totalPullRequestReviewContributionsQuery, varGQLTotalPullRequestReviewContributionsQuery)
	if err != nil {
		return nil, err
	}

	return aladino.BuildIntValue(int(totalPullRequestReviewContributionsQuery.User.ContributionsCollection.TotalPullRequestReviewContributions)), nil
}
