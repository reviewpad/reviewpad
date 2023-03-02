// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Organization() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           organizationCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func organizationCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	orgName := e.GetTarget().GetTargetEntity().Owner
	users, _, err := e.GetGithubClient().ListOrganizationMembers(e.GetCtx(), orgName, nil)
	if err != nil {
		return nil, err
	}
	elems := make([]aladino.Value, len(users))
	for i, user := range users {
		elems[i] = aladino.BuildStringValue(*user.Login)
	}

	return aladino.BuildArrayValue(elems), nil
}
