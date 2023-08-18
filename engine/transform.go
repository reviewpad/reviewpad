// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"fmt"
	"regexp"
	"strings"
)

func addDefaultsToRequestedReviewers(str string) string {
	m := regexp.MustCompile(`\$assignReviewer\(([^,]*\[[^]]*\]|[^,]+)(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?\)`)
	match := m.FindStringSubmatch(str)

	if len(match) == 0 {
		return str
	}

	reviewers := match[1]
	totalRequiredReviewers := `99`
	policy := `"reviewpad"`

	if match[2] != "" {
		totalRequiredReviewers = match[2]
	}

	if match[3] != "" {
		policy = match[3]
	}

	return fmt.Sprintf("$assignReviewer(%s, %s, %s)", reviewers, totalRequiredReviewers, policy)
}

func addDefaultsToAddReviewers(str string) string {
	m := regexp.MustCompile(`\$addReviewers\(([^,]*\[[^]]*\]|[^,]+)(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?\)`)
	match := m.FindStringSubmatch(str)

	if len(match) == 0 {
		return str
	}

	reviewers := match[1]
	totalRequiredReviewers := `99`
	policy := `"reviewpad"`

	if match[2] != "" {
		totalRequiredReviewers = match[2]
	}

	if match[3] != "" {
		policy = match[3]
	}

	return fmt.Sprintf("$addReviewers(%s, %s, %s)", reviewers, totalRequiredReviewers, policy)
}

func addDefaultMergeMethod(str string) string {
	return strings.ReplaceAll(str, "$merge()", "$merge(\"merge\")")
}

func addDefaultSizeMethod(str string) string {
	return strings.ReplaceAll(str, "$size()", "$size([])")
}

func addDefaultGetSizeMethod(str string) string {
	return strings.ReplaceAll(str, "$getSize()", "$getSize([])")
}

func addDefaultIssueCountBy(str string) string {
	r := regexp.MustCompile(`\$issueCountBy\(([^,\s]+)(?:,\s*(("[^\(\)]*")|(\$\w+)))?\)`)
	str = r.ReplaceAllString(str, `$$issueCountBy($1, $2)`)
	r = regexp.MustCompile(`\$issueCountBy\(([^,\s]+),\s*("")?\)`)
	return r.ReplaceAllString(str, `$$issueCountBy($1, "all")`)
}

func addDefaultCountUserIssues(str string) string {
	r := regexp.MustCompile(`\$countUserIssues\(([^,\s]+)(?:,\s*(("[^\(\)]*")|(\$\w+)))?\)`)
	str = r.ReplaceAllString(str, `$$countUserIssues($1, $2)`)
	r = regexp.MustCompile(`\$countUserIssues\(([^,\s]+),\s*("")?\)`)
	return r.ReplaceAllString(str, `$$countUserIssues($1, "all")`)
}

func addDefaultPullRequestCountBy(str string) string {
	r := regexp.MustCompile(`\$pullRequestCountBy\(([^,\s]+)(?:,\s*(("[^\(\)]*")|(\$\w+)))?\)`)
	str = r.ReplaceAllString(str, `$$pullRequestCountBy($1, $2)`)
	r = regexp.MustCompile(`\$pullRequestCountBy\(([^,\s]+),\s*("")?\)`)
	return r.ReplaceAllString(str, `$$pullRequestCountBy($1, "all")`)
}

func addDefaultCountUserPullRequests(str string) string {
	r := regexp.MustCompile(`\$countUserPullRequests\(([^,\s]+)(?:,\s*(("[^\(\)]*")|(\$\w+)))?\)`)
	str = r.ReplaceAllString(str, `$$countUserPullRequests($1, $2)`)
	r = regexp.MustCompile(`\$countUserPullRequests\(([^,\s]+),\s*("")?\)`)
	return r.ReplaceAllString(str, `$$countUserPullRequests($1, "all")`)
}

func addEmptyCloseComment(str string) string {
	return strings.ReplaceAll(str, "$close()", "$close(\"\", \"completed\")")
}

func addDefaultCloseReason(str string) string {
	r := regexp.MustCompile(`\$close\(((?:"([^"]*)")|(?:\$\w+))\)`)
	return r.ReplaceAllString(str, `$$close($1, "completed")`)
}

func addDefaultHaveAllChecksRunCompleted(str string) string {
	allArgsRegex := regexp.MustCompile(`\$haveAllChecksRunCompleted\(([^,]*\[[^]]*\]|[^,]+)?(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?\)`)
	match := allArgsRegex.FindStringSubmatch(str)

	if len(match) == 0 {
		return str
	}

	checkRunsToIgnore := `[]`
	conclusion := `""`
	checkConclusionsToIgnore := `[]`

	if match[1] != "" {
		checkRunsToIgnore = match[1]
	}

	if match[2] != "" {
		conclusion = match[2]
	}

	if match[3] != "" {
		checkConclusionsToIgnore = match[3]
	}

	return fmt.Sprintf(`$haveAllChecksRunCompleted(%s, %s, %s)`, checkRunsToIgnore, conclusion, checkConclusionsToIgnore)
}

func addDefaultHasOnlyCompletedCheckRuns(str string) string {
	allArgsRegex := regexp.MustCompile(`\$hasOnlyCompletedCheckRuns\(([^,]*\[[^]]*\]|[^,]+)?(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?(?:\s*,\s*([^,]*\[[^]]*\]|[^,]+))?\)`)
	match := allArgsRegex.FindStringSubmatch(str)

	if len(match) == 0 {
		return str
	}

	checkRunsToIgnore := `[]`
	conclusion := `""`
	checkConclusionsToIgnore := `[]`

	if match[1] != "" {
		checkRunsToIgnore = match[1]
	}

	if match[2] != "" {
		conclusion = match[2]
	}

	if match[3] != "" {
		checkConclusionsToIgnore = match[3]
	}

	return fmt.Sprintf(`$hasOnlyCompletedCheckRuns(%s, %s, %s)`, checkRunsToIgnore, conclusion, checkConclusionsToIgnore)
}

func addDefaultJoinSeparator(str string) string {
	r := regexp.MustCompile(`\$join\((\[[^\(\)]+\]|(?:[^,]+))(?:,\s*((?:"([^"]*)")|(?:\$\w+)))?\)`)
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
	r = regexp.MustCompile(`\$assignCodeAuthorReviewers\(((?:\d+)|(?:\$\w+)),\s*((?:\[.*\])|(?:[^,]+))\)`)
	return r.ReplaceAllString(str, `$$assignCodeAuthorReviewers($1, $2, 0)`)
}

func addDefaultAddReviewersBasedOnCodeAuthor(str string) string {
	allArgsRegex := regexp.MustCompile(`\$addReviewersBasedOnCodeAuthor\((\d+),\s*((?:\[.*\])|(?:[^,]+)),\s*\d+\)`)
	if allArgsRegex.MatchString(str) {
		return str
	}

	str = strings.ReplaceAll(str, "$addReviewersBasedOnCodeAuthor()", "$addReviewersBasedOnCodeAuthor(1)")
	r := regexp.MustCompile(`\$addReviewersBasedOnCodeAuthor\(([^,]*)\)`)
	str = r.ReplaceAllString(str, `$$addReviewersBasedOnCodeAuthor($1, [])`)
	r = regexp.MustCompile(`\$addReviewersBasedOnCodeAuthor\(((?:\d+)|(?:\$\w+)),\s*((?:\[.*\])|(?:[^,]+))\)`)
	return r.ReplaceAllString(str, `$$addReviewersBasedOnCodeAuthor($1, $2, 0)`)
}

func addDefaultHasAnyCheckRunCompleted(str string) string {
	allArgsRegex := regexp.MustCompile(`\$hasAnyCheckRunCompleted\(((?:\[[^\(\)]*\])|(?:\$.+\(.*\))),\s*((?:\[[^\(\)]*\])|(?:\$.+\(.*\)))\)`)
	if allArgsRegex.MatchString(str) {
		return str
	}

	str = strings.ReplaceAll(str, "$hasAnyCheckRunCompleted()", "$hasAnyCheckRunCompleted([])")

	r := regexp.MustCompile(`\$hasAnyCheckRunCompleted\(((?:\[[^\(\)\[\]]*\])|(?:\$.+\([^\)\(]*\)))\)`)
	return r.ReplaceAllString(str, `$$hasAnyCheckRunCompleted($1, [])`)
}

func addDefaultsToRequestedAssignees(str string) string {
	m := regexp.MustCompile(`\$assignAssignees\(((\[.*\])|(\$(team|group)\("[^"]*"\))|(\$\w+))(,(((-)?(\d+))|(\$\w+)))?\)`)
	match := m.FindStringSubmatch(str)

	if len(match) == 0 {
		return str
	}

	assignees := match[1]
	totalRequiredAssignees := match[6]

	if totalRequiredAssignees == "" {
		totalRequiredAssignees = "10"
	}

	return strings.Replace(str, match[0], "$assignAssignees("+assignees+", "+totalRequiredAssignees+")", 1)
}

func addDefaultsToAddAssignees(str string) string {
	m := regexp.MustCompile(`\$addAssignees\(((\[.*\])|(\$(team|group)\("[^"]*"\))|(\$\w+))(,(((-)?(\d+))|(\$\w+)))?\)`)
	match := m.FindStringSubmatch(str)

	if len(match) == 0 {
		return str
	}

	assignees := match[1]
	totalRequiredAssignees := match[6]

	if totalRequiredAssignees == "" {
		totalRequiredAssignees = "10"
	}

	return strings.Replace(str, match[0], "$addAssignees("+assignees+", "+totalRequiredAssignees+")", 1)
}

func addEmptyFilterToHasCodeWithoutSemanticChanges(str string) string {
	return strings.ReplaceAll(str, "$hasCodeWithoutSemanticChanges()", "$hasCodeWithoutSemanticChanges([])")
}

func addEmptyFilterToContainsOnlyCodeWithoutSemanticChanges(str string) string {
	return strings.ReplaceAll(str, "$containsOnlyCodeWithoutSemanticChanges()", "$containsOnlyCodeWithoutSemanticChanges([])")
}

func summarizeAlias(str string) string {
	return strings.ReplaceAll(str, "$summarize()", `$robinSummarize("default", "openai-gpt-4")`)
}

func addDefaultsGetReviewers(str string) string {
	return strings.ReplaceAll(str, "$getReviewers()", "$getReviewers(\"\")")
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
		addDefaultHasAnyCheckRunCompleted,
		addDefaultsToRequestedAssignees,
		addEmptyFilterToHasCodeWithoutSemanticChanges,
		summarizeAlias,
		addDefaultsGetReviewers,
		addDefaultsToAddReviewers,
		addDefaultGetSizeMethod,
		addDefaultCountUserIssues,
		addDefaultCountUserPullRequests,
		addDefaultHasOnlyCompletedCheckRuns,
		addDefaultAddReviewersBasedOnCodeAuthor,
		addDefaultsToAddAssignees,
		addEmptyFilterToContainsOnlyCodeWithoutSemanticChanges,
	}

	for i := range transformations {
		transformedActionStr = transformations[i](transformedActionStr)
	}

	return transformedActionStr
}
