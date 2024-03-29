// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func FileCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildIntType()),
		Code:           fileCountCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func fileCountCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	return lang.BuildIntValue(len(patch)), nil
}
