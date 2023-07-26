// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func DictionaryValueFromKey() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildStringType()),
		Code:           dictionaryValueFromKeyCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func dictionaryValueFromKeyCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	dictionaryName := args[0].(*lang.StringValue).Val
	dictionaryKeyName := args[1].(*lang.StringValue).Val

	if val, ok := e.GetRegisterMap()[fmt.Sprintf("@dictionary:%s", dictionaryName)]; ok {
		if keyValue, ok := val.(*lang.DictionaryValue).Vals[dictionaryKeyName]; ok {
			return keyValue, nil
		} else {
			return nil, fmt.Errorf("no key with name %v in dictionary %v", dictionaryKeyName, dictionaryName)
		}
	}

	return nil, fmt.Errorf("no dictionary with name %v", dictionaryName)
}
