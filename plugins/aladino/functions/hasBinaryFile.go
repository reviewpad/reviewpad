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

func HasBinaryFile() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildBoolType()),
		Code:           hasBinaryFile,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasBinaryFile(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	target := e.GetTarget().(*target.PullRequestTarget)
	headBranch := target.PullRequest.Head.Name

	for _, patchFile := range target.Patch {
		isBinary, err := target.IsFileBinary(headBranch, patchFile.Repr.GetFilename())
		if err != nil {
			return nil, err
		}

		if isBinary {
			return lang.BuildBoolValue(true), nil
		}
	}

	return lang.BuildBoolValue(false), nil
}
