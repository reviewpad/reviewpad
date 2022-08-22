// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func CreatedAt() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: createdAtCode,
	}
}

func createdAtCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	return aladino.BuildIntValue(int(e.GetPullRequest().GetCreatedAt().Unix())), nil
}
