// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_pullRequest

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

func FileCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: fileCountCode,
	}
}

func fileCountCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	patch := e.GetPatch()
	return aladino.BuildIntValue(len(patch)), nil
}
