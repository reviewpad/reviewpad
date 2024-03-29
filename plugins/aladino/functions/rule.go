// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Rule() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           ruleCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func ruleCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	ruleName := args[0].(*lang.StringValue).Val

	internalRuleName := aladino.BuildInternalRuleName(ruleName)

	if spec, ok := e.GetRegisterMap()[internalRuleName]; ok {
		specRaw := spec.(*lang.StringValue).Val
		result, err := aladino.EvalExpr(e, "patch", specRaw)
		if err != nil {
			return nil, err
		}
		return lang.BuildBoolValue(result), nil
	}

	return nil, fmt.Errorf("$rule: no rule with name %v in state %+q", ruleName, e.GetRegisterMap())
}
