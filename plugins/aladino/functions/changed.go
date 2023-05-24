// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"regexp"
	"strings"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Changed() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           changedCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func changedCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	antecedentRegex := args[0].(*lang.StringValue).Val
	consequentRegex := args[1].(*lang.StringValue).Val

	antecedentMatches := getMatches(pullRequest, antecedentRegex)
	consequentMatches := getMatches(pullRequest, consequentRegex)

	retValue := lang.BuildTrueValue()

	for varName, antecedentVals := range antecedentMatches {
		consequentVals, ok := consequentMatches[varName]
		if !ok {
			return lang.BuildFalseValue(), nil
		}

		for _, leftVal := range antecedentVals {
			leftValExists := false
			for _, rightVal := range consequentVals {
				if leftVal == rightVal {
					leftValExists = true
					break
				}
			}
			if !leftValExists {
				return lang.BuildFalseValue(), nil
			}
		}
	}

	return retValue, nil
}

func getMatches(pullRequest *target.PullRequestTarget, pattern string) map[string][]string {
	resolvedPattern, vars := interpolateRegex(pattern)
	re := regexp.MustCompile(resolvedPattern)

	valsMatrix := make(map[string][]string, 0)

	for fp := range pullRequest.Patch {
		for idx, ranges := range re.FindAllStringSubmatchIndex(fp, -1) {
			lower := ranges[2]
			upper := ranges[3]
			varName := vars[idx]
			val := fp[lower:upper]
			vals, ok := valsMatrix[varName]
			if !ok {
				valsMatrix[varName] = []string{val}
			} else {
				valsMatrix[varName] = append(vals, val)
			}
		}
	}

	return valsMatrix
}

func interpolateRegex(s string) (string, []string) {
	vars := make([]string, 0)
	var sb strings.Builder

	re := regexp.MustCompile(`@[0-9]*\.?[0-9]+`)
	index := 0

	matches := re.FindAllStringIndex(s, -1)
	for _, ranges := range matches {
		lower := ranges[0]
		upper := ranges[1]

		sb.WriteString(s[index:lower])
		index = upper

		varName := s[lower:upper]

		sb.WriteString("(.*)")
		vars = append(vars, varName)

	}

	if index < len(s) {
		sb.WriteString(s[index:])
	}

	return sb.String(), vars
}
