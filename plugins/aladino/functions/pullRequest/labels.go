// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_pullRequest

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

func Labels() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: labelsCode,
	}
}

func labelsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	ghLabels := e.GetPullRequest().Labels
	labels := make([]aladino.Value, len(ghLabels))

	for i, ghLabel := range ghLabels {
		labels[i] = aladino.BuildStringValue(ghLabel.GetName())
	}

	return aladino.BuildArrayValue(labels), nil
}
