// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_organization

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

func Organization() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: organizationCode,
	}
}

func organizationCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	orgName := e.GetPullRequest().Head.Repo.Owner.Login
	users, _, err := e.GetClient().Organizations.ListMembers(e.GetCtx(), *orgName, nil)
	if err != nil {
		return nil, err
	}
	elems := make([]aladino.Value, len(users))
	for i, user := range users {
		elems[i] = aladino.BuildStringValue(*user.Login)
	}

	return aladino.BuildArrayValue(elems), nil
}
