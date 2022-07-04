// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_pullRequest

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

func HasFileName() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: hasFileNameCode,
	}
}

func hasFileNameCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	fileNameStr := args[0].(*aladino.StringValue)

	patch := e.GetPatch()
	for fp := range patch {
		if fp == fileNameStr.Val {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}