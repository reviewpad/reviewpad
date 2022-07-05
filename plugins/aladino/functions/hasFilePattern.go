// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	doublestar "github.com/bmatcuk/doublestar"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
)

func HasFilePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: hasFilePatternCode,
	}
}

func hasFilePatternCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	filePatternRegex := args[0].(*aladino.StringValue)

	patch := e.GetPatch()
	for fp := range patch {
		re, err := doublestar.Match(filePatternRegex.Val, fp)
		if err != nil {
			return aladino.BuildFalseValue(), err
		}
		if re {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}
