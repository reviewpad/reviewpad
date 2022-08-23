// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Group() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           groupCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func groupCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	groupName := args[0].(*aladino.StringValue).Val

	if val, ok := e.GetRegisterMap()[groupName]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("getGroup: no group with name %v in state %+q", groupName, e.GetRegisterMap())
}
