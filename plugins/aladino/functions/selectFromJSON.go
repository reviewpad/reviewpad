// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"encoding/json"
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func SelectFromJSON() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildJSONType(), aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           selectFromContext,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func selectFromJSONCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	jsonValue := args[0].(*aladino.JSONValue).Val
	expr := args[1].(*aladino.StringValue).Val

	parsedExpression, err := jp.ParseString(expr)
	if err != nil {
		return nil, err
	}

	results := parsedExpression.Get(jsonValue)

	if len(results) == 0 {
		return nil, fmt.Errorf(`nothing found at path "%s"`, expr)
	}

	var result interface{} = results

	if len(results) == 1 {
		result = results[0]
	}

	switch res := result.(type) {
	case string:
		return aladino.BuildStringValue(res), nil
	}

	res, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return aladino.BuildStringValue(string(res)), nil
}
