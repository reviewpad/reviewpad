// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func FilesPath() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           filesPathCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func filesPathCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)
	filesPath := make([]lang.Value, 0)

	for _, patchFile := range t.Patch {
		if patchFile.Repr == nil {
			continue
		}

		// TODO: investigate if this ever happens
		// and if we should be using previous file name in this case
		if patchFile.Repr.Filename == "" {
			continue
		}

		filesPath = append(filesPath, lang.BuildStringValue(patchFile.Repr.Filename))
	}

	return lang.BuildArrayValue(filesPath), nil
}
