// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/ohler55/ojg/oj"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ToJSON() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildJSONType()),
		Code:           toJSONCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func toJSONCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	jsonString := args[0].(*aladino.StringValue).Val

	val, err := oj.ParseString(jsonString)
	if err != nil {
		return nil, err
	}

	return aladino.BuildJSONValue(val), nil
}
