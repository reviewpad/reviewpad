// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func HasBinaryFile() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           hasBinaryFile,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasBinaryFile(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	hasBinaryFile := false

	for _, patch := range patch {
		// there won't be any patch for binary files, we can use that to check if pr has binary file changes
		if patch.Repr.GetPatch() == "" {
			hasBinaryFile = true
			break
		}
	}

	return aladino.BuildBoolValue(hasBinaryFile), nil
}
