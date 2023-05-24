// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func SelectFromContext() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           selectFromContext,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func selectFromContext(e aladino.Env, args []lang.Value) (lang.Value, error) {
	path := args[0].(*lang.StringValue)

	targetContext, err := contextCode(e, nil)
	if err != nil {
		return nil, err
	}

	contextJSON, err := toJSONCode(e, []lang.Value{targetContext})
	if err != nil {
		return nil, err
	}

	return selectFromJSONCode(e, []lang.Value{contextJSON, path})
}
