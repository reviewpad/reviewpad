// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_pullRequest

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
)

func IsDraft() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: isDraftCode,
	}
}

func isDraftCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	if pullRequest == nil {
		return nil, fmt.Errorf("isDraft: pull request is nil")
	}

	return aladino.BuildBoolValue(pullRequest.GetDraft()), nil
}
