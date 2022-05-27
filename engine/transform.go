// Copyright (C) 2019-2022 Explore.dev - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

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
