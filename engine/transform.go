// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"regexp"
	"strings"
)

func addDefaultTotalRequestedReviewers(str string) string {
	squadFunctionsRegex := "(\\$(team|group)\\(\"[^\"]*\"\\))"
	reviewersListRegex := "(\\[(.*)\\])"

	m := regexp.MustCompile("\\$assignReviewer\\((" + reviewersListRegex + "|" + squadFunctionsRegex + ")\\)")

	match := m.FindString(str)

	if match == "" {
		return str
	}

	matchWithDefaultNrOfRequestedReviewers := match[0:len(match)-1] + ", 99)"
	transformedStr := strings.ReplaceAll(str, match, matchWithDefaultNrOfRequestedReviewers)

	return transformedStr
}

func addDefaultMergeMethod(str string) string {
	return strings.ReplaceAll(str, "$merge()", "$merge(\"merge\")")
}

func transformActionStr(str string) string {
	transformedActionStr := str

	var transformations = [](func(str string) string){
		addDefaultTotalRequestedReviewers,
		addDefaultMergeMethod,
	}

	for i := range transformations {
		transformedActionStr = transformations[i](transformedActionStr)
	}

	return transformedActionStr
}
