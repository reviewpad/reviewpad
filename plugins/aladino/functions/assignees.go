// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Assignees() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           assigneesCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func assigneesCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	t := e.GetTarget()

	ghAssignees, err := t.GetAssignees()
	if err != nil {
		return nil, err
	}

	assignees := make([]aladino.Value, len(ghAssignees))

	for i, ghAssignee := range ghAssignees {
		assignees[i] = aladino.BuildStringValue(ghAssignee.Login)
	}

	return aladino.BuildArrayValue(assignees), nil
}
