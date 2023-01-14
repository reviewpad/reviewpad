// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func IsLinkedToProject() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           isLinkedToProject,
		SupportedKinds: []handler.TargetEntityKind{handler.Issue, handler.PullRequest},
	}
}

func isLinkedToProject(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	projectTitle := args[0].(*aladino.StringValue).Val

	result, err := e.GetTarget().IsLinkedToProject(projectTitle)
	if err != nil {
		return nil, err
	}

	return aladino.BuildBoolValue(result), nil
}
