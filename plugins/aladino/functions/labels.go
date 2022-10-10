// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Labels() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           labelsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func labelsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	ghLabels := e.GetTarget().GetLabels()

	labels := make([]aladino.Value, len(ghLabels))

	for i, ghLabel := range ghLabels {
		labels[i] = aladino.BuildStringValue(ghLabel.Name)
	}

	return aladino.BuildArrayValue(labels), nil
}
