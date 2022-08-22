// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Comments() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: commentsCode,
	}
}

func commentsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	ghComments, err := e.GetTarget().GetComments()
	if err != nil {
		return nil, err
	}

	commentsBody := make([]aladino.Value, len(ghComments))
	for i, ghComment := range ghComments {
		commentsBody[i] = aladino.BuildStringValue(ghComment.Body)
	}

	return aladino.BuildArrayValue(commentsBody), nil
}
