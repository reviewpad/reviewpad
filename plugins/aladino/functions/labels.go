// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Labels() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           labelsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func labelsCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	ghLabels := e.GetTarget().GetLabels()

	labels := make([]lang.Value, len(ghLabels))

	for i, ghLabel := range ghLabels {
		labels[i] = lang.BuildStringValue(ghLabel.Name)
	}

	return lang.BuildArrayValue(labels), nil
}
