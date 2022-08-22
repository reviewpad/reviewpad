// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func HasLinkedIssues() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasLinkedIssuesCode,
	}
}

func hasLinkedIssuesCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	closingIssuesCount, err := e.GetTarget().GetLinkedIssuesCount()
	if err != nil {
		return nil, err
	}

	return aladino.BuildBoolValue(closingIssuesCount > 0), nil
}
