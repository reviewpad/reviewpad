// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"

	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Rule() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           ruleCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func ruleCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	ruleName := args[0].(*aladino.StringValue).Val

	internalRuleName := aladino.BuildInternalRuleName(ruleName)

	if spec, ok := e.GetRegisterMap()[internalRuleName]; ok {
		specRaw := spec.(*aladino.StringValue).Val
		result, err := aladino.EvalExpr(e, "patch", specRaw)
		if err != nil {
			return nil, err
		}
		return aladino.BuildBoolValue(result), nil
	}

	return nil, fmt.Errorf("$rule: no rule with name %v in state %+q", ruleName, e.GetRegisterMap())
}
