// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/shurcooL/githubv4"
)

func TotalCodeReviews() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildIntType()),
		Code:           totalCodeReviewsCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func totalCodeReviewsCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	username := args[0].(*lang.StringValue).Val

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

	return lang.BuildIntValue(int(totalPullRequestReviewContributionsQuery.User.ContributionsCollection.TotalPullRequestReviewContributions)), nil
}
