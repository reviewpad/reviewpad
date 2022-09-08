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

func addDefaultSizeMethod(str string) string {
	return strings.ReplaceAll(str, "$size()", "$size([])")
}

// reviewpad-an: generated-by-co-pilot
func addDefaultIssueCountBy(str string) string {
	m := regexp.MustCompile("\\$issueCountBy\\(\"[^\"]*\"\\)")
	match := m.FindString(str)

	if match == "" {
		return str
	}

	matchWithDefaultState := match[0:len(match)-1] + ", \"all\")"
	transformedStr := strings.ReplaceAll(str, match, matchWithDefaultState)
	return transformedStr
}

func addDefaultPullRequestCountBy(str string) string {
	m := regexp.MustCompile("\\$pullRequestCountBy\\(\"[^\"]*\"\\)")
	match := m.FindString(str)

	if match == "" {
		return str
	}

	matchWithDefaultState := match[0:len(match)-1] + ", \"all\")"
	transformedStr := strings.ReplaceAll(str, match, matchWithDefaultState)
	return transformedStr
}

func addEmptyCloseComment(str string) string {
	return strings.ReplaceAll(str, "$close()", "$close(\"\")")
}

func transformAladinoExpression(str string) string {
	transformedActionStr := str

	var transformations = [](func(str string) string){
		addDefaultTotalRequestedReviewers,
		addDefaultMergeMethod,
		addDefaultSizeMethod,
		addDefaultIssueCountBy,
		addDefaultPullRequestCountBy,
		addEmptyCloseComment,
	}

	for i := range transformations {
		transformedActionStr = transformations[i](transformedActionStr)
	}

	return transformedActionStr
}
