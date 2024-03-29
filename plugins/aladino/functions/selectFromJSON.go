// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"encoding/json"

	"github.com/ohler55/ojg/jp"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func SelectFromJSON() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildJSONType(), lang.BuildStringType()}, lang.BuildStringType()),
		Code:           selectFromJSONCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func selectFromJSONCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	jsonValue := args[0].(*lang.JSONValue).Val
	expr := args[1].(*lang.StringValue).Val
	log := e.GetLogger().WithField("builtin", "selectFromJSON")

	parsedExpression, err := jp.ParseString(expr)
	if err != nil {
		return nil, err
	}

	results := parsedExpression.Get(jsonValue)

	if len(results) == 0 {
		log.Infof(`nothing found at path "%s"\n`, expr)
		return lang.BuildStringValue(""), nil
	}

	var result interface{} = results

	if len(results) == 1 {
		result = results[0]
	}

	// marshaling a string into json will cause it to have quotation around it
	if res, ok := result.(string); ok {
		return lang.BuildStringValue(res), nil
	}

	res, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return lang.BuildStringValue(string(res)), nil
}
