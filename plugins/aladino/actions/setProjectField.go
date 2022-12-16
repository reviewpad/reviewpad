// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func SetProjectField() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType(), aladino.BuildStringType()}, nil),
		Code:           setProjectField,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func setProjectField(e aladino.Env, args []aladino.Value) error {
	projectTitle := args[0].(*aladino.StringValue).Val
	fieldName := args[1].(*aladino.StringValue).Val
	fieldValue := args[2].(*aladino.StringValue).Val
	target := e.GetTarget()

	projectItems, err := target.GetLinkedProjects()
	if err != nil {
		return err
	}

	return target.SetProjectFieldSingleSelect(projectItems, projectTitle, fieldName, fieldValue)
}
