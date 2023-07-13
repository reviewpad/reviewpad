// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Assignees() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           assigneesCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func assigneesCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	t := e.GetTarget()

	ghAssignees := t.GetAssignees()

	assignees := make([]lang.Value, len(ghAssignees))

	for i, ghAssignee := range ghAssignees {
		assignees[i] = lang.BuildStringValue(ghAssignee.Login)
	}

	return lang.BuildArrayValue(assignees), nil
}
