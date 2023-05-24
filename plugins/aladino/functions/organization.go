// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Organization() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           organizationCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func organizationCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	orgName := e.GetTarget().GetTargetEntity().Owner
	users, _, err := e.GetGithubClient().ListOrganizationMembers(e.GetCtx(), orgName, nil)
	if err != nil {
		return nil, err
	}
	elems := make([]lang.Value, len(users))
	for i, user := range users {
		elems[i] = lang.BuildStringValue(*user.Login)
	}

	return lang.BuildArrayValue(elems), nil
}
