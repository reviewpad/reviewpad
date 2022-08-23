// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func AppendString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           appendStringCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func appendStringCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	slice1 := args[0].(*aladino.ArrayValue).Vals
	slice2 := args[1].(*aladino.ArrayValue).Vals

	return aladino.BuildArrayValue(append(slice1, slice2...)), nil
}
