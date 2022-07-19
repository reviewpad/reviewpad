// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import "github.com/reviewpad/reviewpad/v3/lang/aladino"

func Author() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: authorCode,
	}
}

func authorCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	authorLogin := e.GetPullRequest().User.GetLogin()
	return aladino.BuildStringValue(authorLogin), nil
}
