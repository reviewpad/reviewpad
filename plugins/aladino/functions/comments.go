// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Comments() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           commentsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func commentsCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	ghComments, err := e.GetTarget().GetComments()
	if err != nil {
		return nil, err
	}

	commentsBody := make([]lang.Value, len(ghComments))
	for i, ghComment := range ghComments {
		commentsBody[i] = lang.BuildStringValue(ghComment.Body)
	}

	return lang.BuildArrayValue(commentsBody), nil
}
