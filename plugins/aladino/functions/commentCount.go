// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

func CommentCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: commentCountCode,
	}
}

func commentCountCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildIntValue(*e.GetPullRequest().Comments), nil
}
