// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func SetProjectField() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType(), lang.BuildStringType()}, nil),
		Code:           setProjectField,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func setProjectField(e aladino.Env, args []lang.Value) error {
	projectTitle := args[0].(*lang.StringValue).Val
	fieldName := args[1].(*lang.StringValue).Val
	fieldValue := args[2].(*lang.StringValue).Val
	target := e.GetTarget()

	isInProject, err := target.IsInProject(projectTitle)
	if err != nil {
		return err
	}

	if !isInProject {
		err = addToProjectCode(e, []lang.Value{args[0], &lang.StringValue{Val: ""}})
		if err != nil {
			return err
		}
	}

	return target.SetProjectField(projectTitle, fieldName, fieldValue)
}
