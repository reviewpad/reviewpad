// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

func Size() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: sizeCode,
	}
}

func sizeCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	size := e.GetPullRequest().GetAdditions() + e.GetPullRequest().GetDeletions()
	return aladino.BuildIntValue(size), nil
}
