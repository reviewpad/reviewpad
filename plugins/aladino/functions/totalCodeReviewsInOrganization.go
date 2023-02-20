// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func TotalCodeReviewsInOrganization() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildIntType()),
		Code:           totalCodeReviewsInOrganizationCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func totalCodeReviewsInOrganizationCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	username := args[0].(*aladino.StringValue).Val
	entity := e.GetTarget().GetTargetEntity()
	var reviewsCount int

	ghPullRequest, _, fetchPrErr := e.GetGithubClient().GetPullRequest(e.GetCtx(), entity.Owner, entity.Repo, entity.Number)
	if fetchPrErr != nil {
		return nil, fetchPrErr
	}

	// Fetch the repositories owned by the organization/user that owns the repository where the pull request was opened
	repositories, fetchReposErr := e.GetGithubClient().GetRepositoriesByOrgOrUserLogin(e.GetCtx(), host.GetPullRequestHeadOwnerName(ghPullRequest))
	if fetchReposErr != nil {
		return nil, fetchReposErr
	}

	for _, repository := range repositories {
		reviewsCountInRepo, err := e.GetGithubClient().GetReviewsCountByUserFromOpenPullRequests(e.GetCtx(), repository, username)
		if err != nil {
			return nil, err
		}

		reviewsCount += reviewsCountInRepo
	}

	e.GetLogger().Logger.Info(fmt.Sprintf("The user %s has created %d reviews in open pull requests.", username, reviewsCount))

	return aladino.BuildIntValue(reviewsCount), nil
}
