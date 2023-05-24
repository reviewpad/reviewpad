// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IsLinkedToProject() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           isLinkedToProject,
		SupportedKinds: []entities.TargetEntityKind{entities.Issue, entities.PullRequest},
	}
}

func isLinkedToProject(e aladino.Env, args []lang.Value) (lang.Value, error) {
	projectTitle := args[0].(*lang.StringValue).Val

	result, err := e.GetTarget().IsLinkedToProject(projectTitle)
	if err != nil {
		return nil, err
	}

	return lang.BuildBoolValue(result), nil
}
