// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasBinaryFile() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           hasBinaryFile,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasBinaryFile(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	target := e.GetTarget().(*target.PullRequestTarget)
	headBranch := target.PullRequest.GetHeadBranch()

	for _, patchFile := range target.Patch {
		isBinary, err := target.IsFileBinary(headBranch, patchFile.Repr.GetFilename())
		if err != nil {
			return nil, err
		}

		if isBinary {
			return aladino.BuildBoolValue(true), nil
		}
	}

	return aladino.BuildBoolValue(false), nil
}
