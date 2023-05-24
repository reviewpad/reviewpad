// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Group() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           groupCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func groupCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	groupName := args[0].(*lang.StringValue).Val

	if val, ok := e.GetRegisterMap()[groupName]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("getGroup: no group with name %v in state %+q", groupName, e.GetRegisterMap())
}
