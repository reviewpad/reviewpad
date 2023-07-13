// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AppendString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildArrayOfType(lang.BuildStringType()), lang.BuildArrayOfType(lang.BuildStringType())}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           appendStringCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func appendStringCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	slice1 := args[0].(*lang.ArrayValue).Vals
	slice2 := args[1].(*lang.ArrayValue).Vals

	return lang.BuildArrayValue(append(slice1, slice2...)), nil
}
