// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasLabel() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           hasLabelCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func hasLabelCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	label := args[0].(*lang.StringValue)

	labels := e.GetTarget().GetLabels()
	for _, l := range labels {
		if l.Name == label.Val {
			return lang.BuildTrueValue(), nil
		}
	}

	return lang.BuildFalseValue(), nil
}