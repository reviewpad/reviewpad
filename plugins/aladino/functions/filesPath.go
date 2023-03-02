// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func FilesPath() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           filesPathCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func filesPathCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)
	filesPath := make([]aladino.Value, 0)

	for _, patchFile := range t.Patch {
		if patchFile.Repr == nil {
			continue
		}

		if patchFile.Repr.Filename == nil {
			continue
		}

		filesPath = append(filesPath, aladino.BuildStringValue(*patchFile.Repr.Filename))
	}

	return aladino.BuildArrayValue(filesPath), nil
}
