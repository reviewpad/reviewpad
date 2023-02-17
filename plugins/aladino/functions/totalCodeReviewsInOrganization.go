// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"log"
	"strings"

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
		repositoryOwner := strings.Split(repository, "/")[0]
		repositoryName := strings.Split(repository, "/")[1]

		userReviewsCountInRepo, err := e.GetGithubClient().GetReviewsCountByUserFromOpenPullRequests(e.GetCtx(), repositoryOwner, repositoryName, username)
		if err != nil {
			return nil, err
		}

		reviewsCount += userReviewsCountInRepo
	}

	log.Printf("The user %s has created %d reviews in open pull requests.", username, reviewsCount)

	return aladino.BuildIntValue(reviewsCount), nil
}
