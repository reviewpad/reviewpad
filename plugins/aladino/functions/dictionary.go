// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Dictionary() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildDictionaryType()),
		Code:           dictionaryCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func dictionaryCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	dictionaryName := args[0].(*lang.StringValue).Val

	if val, ok := e.GetRegisterMap()[fmt.Sprintf("@dictionary:%s", dictionaryName)]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("getDictionary: no dictionary with name %v in state %+q", dictionaryName, e.GetRegisterMap())
}
