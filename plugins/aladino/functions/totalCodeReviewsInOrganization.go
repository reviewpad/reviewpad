// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
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

	totalReviews, err := e.GetGithubClient().GetTotalReviewsByUserInOrg(e.GetCtx(), e.GetTarget().GetTargetEntity().Owner, username)

	return aladino.BuildIntValue(totalReviews), err
}
