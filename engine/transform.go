// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"regexp"
	"strings"
)

func addDefaultsToRequestedReviewers(str string) string {
	trimmedStr := strings.ReplaceAll(str, " ", "")
	m := regexp.MustCompile(`\$assignReviewer\(((\[.*\])|(\$(team|group)\("[^"]*"\)))(,(\d+))?(,"([^"]*)")?\)`)
	match := m.FindStringSubmatch(trimmedStr)

	if len(match) == 0 {
		return str
	}

	reviewers := match[1]
	totalRequiredReviewers := match[6]
	policy := match[8]

	if totalRequiredReviewers == "" {
		totalRequiredReviewers = "99"
	}

	if policy == "" {
		policy = "reviewpad"
	}

	return strings.Replace(trimmedStr, match[0], "$assignReviewer("+reviewers+", "+totalRequiredReviewers+", \""+policy+"\")", 1)
}

func addDefaultMergeMethod(str string) string {
	return strings.ReplaceAll(str, "$merge()", "$merge(\"merge\")")
}

func addDefaultSizeMethod(str string) string {
	return strings.ReplaceAll(str, "$size()", "$size([])")
}

func addDefaultIssueCountBy(str string) string {
	r := regexp.MustCompile(`\$issueCountBy\(([^,\s]+)(?:,\s*("[^\(\)]*"))?\)`)
	str = r.ReplaceAllString(str, `$$issueCountBy($1, $2)`)
	r = regexp.MustCompile(`\$issueCountBy\(([^,\s]+),\s*("")?\)`)
	return r.ReplaceAllString(str, `$$issueCountBy($1, "all")`)
}

func addDefaultPullRequestCountBy(str string) string {
	r := regexp.MustCompile(`\$pullRequestCountBy\(([^,\s]+)(?:,\s*("[^\(\)]*"))?\)`)
	str = r.ReplaceAllString(str, `$$pullRequestCountBy($1, $2)`)
	r = regexp.MustCompile(`\$pullRequestCountBy\(([^,\s]+),\s*("")?\)`)
	return r.ReplaceAllString(str, `$$pullRequestCountBy($1, "all")`)
}

func addEmptyCloseComment(str string) string {
	return strings.ReplaceAll(str, "$close()", "$close(\"\", \"completed\")")
}

func addDefaultCloseReason(str string) string {
	r := regexp.MustCompile(`\$close\("([^"]*)"\)`)
	return r.ReplaceAllString(str, `$$close("$1", "completed")`)
}

func addDefaultHaveAllChecksRunCompleted(str string) string {
	r := regexp.MustCompile(`\$haveAllChecksRunCompleted\(((?:\[(?:".*")?(?:,".*")?\])|(?:[^,\s\[\]]+))?(?:,\s*("[^\(\)]*"))?\)`)
	str = r.ReplaceAllString(str, `$$haveAllChecksRunCompleted($1, $2)`)
	r = regexp.MustCompile(`\$haveAllChecksRunCompleted\(((?:\[(?:".*")?(?:,".*")?\])|(?:[^,\s]+))?(\,\s?)?\)`)
	str = r.ReplaceAllString(str, `$$haveAllChecksRunCompleted($1, "")`)
	r = regexp.MustCompile(`\$haveAllChecksRunCompleted\(,\s?""\)`)
	return r.ReplaceAllString(str, `$$haveAllChecksRunCompleted([], "")`)
}

func addDefaultJoinSeparator(str string) string {
	r := regexp.MustCompile(`\$join\((\[[^\(\)]+\]|(?:[^,]+))(?:,\s*("([^"]*)"))?\)`)
	str = r.ReplaceAllString(str, `$$join($1, $2)`)
	r = regexp.MustCompile(`\$join\((\[.*\]|(?:[^,]*)), \)`)
	return r.ReplaceAllString(str, `$$join($1, " ")`)
}

func addEmptyApproveComment(str string) string {
	return strings.ReplaceAll(str, "$approve()", "$approve(\"\")")
}

func addDefaultAssignCodeAuthorReviewer(str string) string {
	allArgsRegex := regexp.MustCompile(`\$assignCodeAuthorReviewers\((\d+),\s*((?:\[.*\])|(?:[^,]+)),\s*\d+\)`)
	if allArgsRegex.MatchString(str) {
		return str
	}

	str = strings.ReplaceAll(str, "$assignCodeAuthorReviewers()", "$assignCodeAuthorReviewers(1)")
	r := regexp.MustCompile(`\$assignCodeAuthorReviewers\(([^,]*)\)`)
	str = r.ReplaceAllString(str, `$$assignCodeAuthorReviewers($1, [])`)
	r = regexp.MustCompile(`\$assignCodeAuthorReviewers\((\d+),\s*((?:\[.*\])|(?:[^,]+))\)`)
	return r.ReplaceAllString(str, `$$assignCodeAuthorReviewers($1, $2, 0)`)
}

func transformAladinoExpression(str string) string {
	transformedActionStr := str

	var transformations = [](func(str string) string){
		addDefaultsToRequestedReviewers,
		addDefaultMergeMethod,
		addDefaultSizeMethod,
		addDefaultIssueCountBy,
		addDefaultPullRequestCountBy,
		addEmptyCloseComment,
		addDefaultCloseReason,
		addDefaultHaveAllChecksRunCompleted,
		addDefaultJoinSeparator,
		addEmptyApproveComment,
		addDefaultAssignCodeAuthorReviewer,
	}

	for i := range transformations {
		transformedActionStr = transformations[i](transformedActionStr)
	}

	return transformedActionStr
}
