// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_utils

import (
	"strings"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
)

func Contains() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: containsCode,
	}
}

func containsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	str := args[0].(*aladino.StringValue).Val
	subString := args[1].(*aladino.StringValue).Val

	return aladino.BuildBoolValue(strings.Contains(str, subString)), nil
}
