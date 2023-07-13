// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"strconv"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ToBool() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           toBoolCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func toBoolCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	str := args[0].(*lang.StringValue).Val

	val, err := strconv.ParseBool(str)
	if err != nil {
		return nil, fmt.Errorf(`error converting "%s" to boolean: %w`, str, err)
	}

	return lang.BuildBoolValue(val), nil
}
