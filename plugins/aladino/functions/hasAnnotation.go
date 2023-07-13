// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"strings"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	semantic "github.com/reviewpad/reviewpad/v4/plugins/aladino/semantic"
)

func HasAnnotation() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildBoolType()),
		Code:           hasAnnotationCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasAnnotationCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	annotation := args[0].(*lang.StringValue).Val

	symbolsByFileName, err := semantic.GetSymbolsFromPatch(e)
	if err != nil {
		return nil, err
	}

	for _, fileSymbols := range symbolsByFileName {
		for _, symbol := range fileSymbols.Symbols {
			for _, symbolComment := range symbol.CodeComments {
				comment := symbolComment.Code
				if commentHasAnnotation(comment, annotation) {
					return lang.BuildTrueValue(), nil
				}
			}
		}
	}

	return lang.BuildFalseValue(), nil
}

func commentHasAnnotation(comment, annotation string) bool {
	normalizedComment := strings.ToLower(strings.Trim(comment, " \n"))
	normalizedAnnotation := strings.ToLower(annotation)
	anPrefix := "reviewpad-an: "
	if strings.HasPrefix(normalizedComment, anPrefix) {
		rest := normalizedComment[len(anPrefix)-1:]
		restUntilNewLine := strings.Split(rest, "\n")[0]

		annotations := strings.Fields(restUntilNewLine)
		for _, annot := range annotations {
			if normalizedAnnotation == annot {
				return true
			}
		}
	}

	return false
}
