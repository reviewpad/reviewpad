// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	semantic "github.com/reviewpad/reviewpad/v4/plugins/aladino/semantic"
)

// RequestIDKey identifies request id field in context
const RequestIDKey = "request-id"

func HasCodeWithoutSemanticChanges() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildArrayOfType(lang.BuildStringType())}, lang.BuildBoolType()),
		Code:           hasCodeWithoutSemanticChanges,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasCodeWithoutSemanticChanges(e aladino.Env, args []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	log := e.GetLogger()

	// filter patch by file pattern
	patch := pullRequest.Patch
	newPatch := make(map[string]*codehost.File)

	filePatterns := args[0].(*lang.ArrayValue)
	for fp, file := range patch {
		ignore := false
		for _, v := range filePatterns.Vals {
			filePatternRegex := v.(*lang.StringValue)

			re, err := doublestar.Match(filePatternRegex.Val, fp)
			if err == nil && re {
				ignore = true
				break
			}
		}
		if !ignore {
			newPatch[fp] = file
			log.Infof("file %s is not ignored", fp)
		} else {
			log.Infof("file %s is ignored", fp)
		}
	}

	baseSymbols, _, err := semantic.GetSymbolsFromBaseByPatch(e, newPatch)
	if err != nil {
		return nil, err
	}

	headSymbols, _, err := semantic.GetSymbolsFromHeadByPatch(e, newPatch)
	if err != nil {
		return nil, err
	}

	// match the symbols based on their ID
	res := true

	if len(baseSymbols.Symbols) != len(headSymbols.Symbols) {
		res = false
	}

	for symbolID, symbol := range baseSymbols.Symbols {
		_, ok := headSymbols.Symbols[symbolID]
		if !ok {
			log.Infof("symbol %v in %v is different", symbol.Name, symbol.Definition.FilePath)
			res = false
			break
		}
	}

	if res {
		log.Infof("no semantic changes detected")
		return lang.BuildTrueValue(), nil
	} else {
		log.Infof("semantic changes detected")
		return lang.BuildFalseValue(), nil
	}
}
