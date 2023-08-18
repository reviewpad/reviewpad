// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformAladinoExpression(t *testing.T) {
	tests := map[string]struct {
		arg     string
		wantVal string
	}{
		"id": {
			arg:     "$id()",
			wantVal: "$id()",
		},
		"merge": {
			arg:     "$merge()",
			wantVal: "$merge(\"merge\")",
		},
		"merge with variable": {
			arg:     "$merge($mergeMethod)",
			wantVal: "$merge($mergeMethod)",
		},
		"size": {
			arg:     "$size()",
			wantVal: "$size([])",
		},
		"size with variable": {
			arg:     "$size($files)",
			wantVal: "$size($files)",
		},
		"get size": {
			arg:     "$getSize()",
			wantVal: "$getSize([])",
		},
		"get size with variable": {
			arg:     "$getSize($files)",
			wantVal: "$getSize($files)",
		},
		"issueCountBy simple": {
			arg:     "$issueCountBy(\"john\", \"open\") > 0",
			wantVal: "$issueCountBy(\"john\", \"open\") > 0",
		},
		"issueCountBy": {
			arg:     "$issueCountBy(\"john\", \"open\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$issueCountBy(\"john\", \"open\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"issueCountBy with one variable": {
			arg:     `$issueCountBy($teamMember)`,
			wantVal: `$issueCountBy($teamMember, "all")`,
		},
		"issueCountBy with two variables": {
			arg:     `$issueCountBy($teamMember, $state)`,
			wantVal: `$issueCountBy($teamMember, $state)`,
		},
		"issueCountBy with last one variable": {
			arg:     `$issueCountBy($author(), $state)`,
			wantVal: `$issueCountBy($author(), $state)`,
		},
		"countUserIssues simple": {
			arg:     "$countUserIssues(\"john\", \"open\") > 0",
			wantVal: "$countUserIssues(\"john\", \"open\") > 0",
		},
		"countUserIssues": {
			arg:     "$countUserIssues(\"john\", \"open\") > 0 && true && $countUserIssues(\"dev\") > 0",
			wantVal: "$countUserIssues(\"john\", \"open\") > 0 && true && $countUserIssues(\"dev\", \"all\") > 0",
		},
		"countUserIssues with one variable": {
			arg:     `$countUserIssues($teamMember)`,
			wantVal: `$countUserIssues($teamMember, "all")`,
		},
		"countUserIssues with two variables": {
			arg:     `$countUserIssues($teamMember, $state)`,
			wantVal: `$countUserIssues($teamMember, $state)`,
		},
		"countUserIssues with last one variable": {
			arg:     `$countUserIssues($author(), $state)`,
			wantVal: `$countUserIssues($author(), $state)`,
		},
		"pullRequestCountBy simple": {
			arg:     "$pullRequestCountBy(\"john\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0",
		},
		"pullRequestCountBy nil state": {
			arg:     "$pullRequestCountBy(\"john\", \"\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0",
		},
		"pullRequestCountBy nil dev": {
			arg:     "$pullRequestCountBy(\"\", \"closed\") > 0",
			wantVal: "$pullRequestCountBy(\"\", \"closed\") > 0",
		},
		"pullRequestCountBy nil dev and state": {
			arg:     "$pullRequestCountBy(\"\") > 0",
			wantVal: "$pullRequestCountBy(\"\", \"all\") > 0",
		},
		"pullRequestCountBy and issueCountBy": {
			arg:     "$pullRequestCountBy(\"john\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"pullRequestCountBy with dev": {
			arg:     `$pullRequestCountBy("john") > 0 && true && $pullRequestCountBy("john", "open") > 1`,
			wantVal: `$pullRequestCountBy("john", "all") > 0 && true && $pullRequestCountBy("john", "open") > 1`,
		},
		"pullRequestCountBy with function call": {
			arg:     `$pullRequestCountBy($author()) > 1 && $pullRequestCountBy($author(), "open") > 2`,
			wantVal: `$pullRequestCountBy($author(), "all") > 1 && $pullRequestCountBy($author(), "open") > 2`,
		},
		"pullRequestCountBy with function call and state": {
			arg:     `$pullRequestCountBy($author(), "open")`,
			wantVal: `$pullRequestCountBy($author(), "open")`,
		},
		"pullRequestCountBy with function call and empty state": {
			arg:     `$pullRequestCountBy($author(), "")`,
			wantVal: `$pullRequestCountBy($author(), "all")`,
		},
		"pullRequestCountBy with one variable": {
			arg:     `$pullRequestCountBy($teamMember)`,
			wantVal: `$pullRequestCountBy($teamMember, "all")`,
		},
		"pullRequestCountBy with two variables": {
			arg:     `$pullRequestCountBy($teamMember, $state)`,
			wantVal: `$pullRequestCountBy($teamMember, $state)`,
		},
		"pullRequestCountBy with last one variable": {
			arg:     `$pullRequestCountBy($author(), $state)`,
			wantVal: `$pullRequestCountBy($author(), $state)`,
		},
		"countUserPullRequests simple": {
			arg:     "$pullRequestCountBy(\"john\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0",
		},
		"countUserPullRequests nil state": {
			arg:     "$countUserPullRequests(\"john\", \"\") > 0",
			wantVal: "$countUserPullRequests(\"john\", \"all\") > 0",
		},
		"countUserPullRequests nil dev": {
			arg:     "$countUserPullRequests(\"\", \"closed\") > 0",
			wantVal: "$countUserPullRequests(\"\", \"closed\") > 0",
		},
		"countUserPullRequests nil dev and state": {
			arg:     "$countUserPullRequests(\"\") > 0",
			wantVal: "$countUserPullRequests(\"\", \"all\") > 0",
		},
		"countUserPullRequests and issueCountBy": {
			arg:     "$countUserPullRequests(\"john\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$countUserPullRequests(\"john\", \"all\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"countUserPullRequests with dev": {
			arg:     `$countUserPullRequests("john") > 0 && true && $countUserPullRequests("john", "open") > 1`,
			wantVal: `$countUserPullRequests("john", "all") > 0 && true && $countUserPullRequests("john", "open") > 1`,
		},
		"countUserPullRequests with function call": {
			arg:     `$countUserPullRequests($author()) > 1 && $countUserPullRequests($author(), "open") > 2`,
			wantVal: `$countUserPullRequests($author(), "all") > 1 && $countUserPullRequests($author(), "open") > 2`,
		},
		"countUserPullRequests with function call and state": {
			arg:     `$countUserPullRequests($author(), "open")`,
			wantVal: `$countUserPullRequests($author(), "open")`,
		},
		"countUserPullRequests with function call and empty state": {
			arg:     `$countUserPullRequests($author(), "")`,
			wantVal: `$countUserPullRequests($author(), "all")`,
		},
		"countUserPullRequests with one variable": {
			arg:     `$countUserPullRequests($teamMember)`,
			wantVal: `$countUserPullRequests($teamMember, "all")`,
		},
		"countUserPullRequests with two variables": {
			arg:     `$countUserPullRequests($teamMember, $state)`,
			wantVal: `$countUserPullRequests($teamMember, $state)`,
		},
		"countUserPullRequests with last one variable": {
			arg:     `$countUserPullRequests($author(), $state)`,
			wantVal: `$countUserPullRequests($author(), $state)`,
		},
		"close": {
			arg:     "$close()",
			wantVal: `$close("", "completed")`,
		},
		"close with comment": {
			arg:     `$close("comment")`,
			wantVal: `$close("comment", "completed")`,
		},
		"close with comment and not_planned state reason": {
			arg:     `$close("comment", "not_planned")`,
			wantVal: `$close("comment", "not_planned")`,
		},
		"close with comment and completed state reason": {
			arg:     `$close("comment", "completed")`,
			wantVal: `$close("comment", "completed")`,
		},
		"close with empty comment and completed state reason": {
			arg:     `$close("", "completed")`,
			wantVal: `$close("", "completed")`,
		},
		"close with empty comment and not_planned state reason": {
			arg:     `$close("", "not_planned")`,
			wantVal: `$close("", "not_planned")`,
		},
		"close with one variable": {
			arg:     "$close($closeComment)",
			wantVal: `$close($closeComment, "completed")`,
		},
		"close with two variable": {
			arg:     "$close($closeComment, $closeState)",
			wantVal: `$close($closeComment, $closeState)`,
		},
		"close with last one variable": {
			arg:     `$close("closing as completed", $closeState)`,
			wantVal: `$close("closing as completed", $closeState)`,
		},
		"$haveAllChecksRunCompleted with arg provided": {
			arg:     `$haveAllChecksRunCompleted(["build", "test"], "", [])`,
			wantVal: `$haveAllChecksRunCompleted(["build", "test"], "", [])`,
		},
		"$haveAllChecksRunCompleted with no arg provided": {
			arg:     `$haveAllChecksRunCompleted()`,
			wantVal: `$haveAllChecksRunCompleted([], "", [])`,
		},
		"$haveAllChecksRunCompleted with no empty arg provided": {
			arg:     `$haveAllChecksRunCompleted(["build", "test"], "success", ["skipped"])`,
			wantVal: `$haveAllChecksRunCompleted(["build", "test"], "success", ["skipped"])`,
		},
		"$haveAllChecksRunCompleted with two args provided": {
			arg:     `$haveAllChecksRunCompleted(["build"], "completed")`,
			wantVal: `$haveAllChecksRunCompleted(["build"], "completed", [])`,
		},
		"$haveAllChecksRunCompleted with one variable": {
			arg:     `$haveAllChecksRunCompleted($checks)`,
			wantVal: `$haveAllChecksRunCompleted($checks, "", [])`,
		},
		"$haveAllChecksRunCompleted with two variables": {
			arg:     `$haveAllChecksRunCompleted($checks, $conclusion)`,
			wantVal: `$haveAllChecksRunCompleted($checks, $conclusion, [])`,
		},
		"$haveAllChecksRunCompleted with three variables": {
			arg:     `$haveAllChecksRunCompleted($checks, $conclusion, $ignoredConclusions)`,
			wantVal: `$haveAllChecksRunCompleted($checks, $conclusion, $ignoredConclusions)`,
		},
		"$hasOnlyCompletedCheckRuns with arg provided": {
			arg:     `$hasOnlyCompletedCheckRuns(["build", "test"], "", [])`,
			wantVal: `$hasOnlyCompletedCheckRuns(["build", "test"], "", [])`,
		},
		"$hasOnlyCompletedCheckRuns with no arg provided": {
			arg:     `$hasOnlyCompletedCheckRuns()`,
			wantVal: `$hasOnlyCompletedCheckRuns([], "", [])`,
		},
		"$hasOnlyCompletedCheckRuns with no empty arg provided": {
			arg:     `$hasOnlyCompletedCheckRuns(["build", "test"], "success", ["skipped"])`,
			wantVal: `$hasOnlyCompletedCheckRuns(["build", "test"], "success", ["skipped"])`,
		},
		"$hasOnlyCompletedCheckRuns with two args provided": {
			arg:     `$hasOnlyCompletedCheckRuns(["build"], "completed")`,
			wantVal: `$hasOnlyCompletedCheckRuns(["build"], "completed", [])`,
		},
		"$hasOnlyCompletedCheckRuns with one variable": {
			arg:     `$hasOnlyCompletedCheckRuns($checks)`,
			wantVal: `$hasOnlyCompletedCheckRuns($checks, "", [])`,
		},
		"$hasOnlyCompletedCheckRuns with two variables": {
			arg:     `$hasOnlyCompletedCheckRuns($checks, $conclusion)`,
			wantVal: `$hasOnlyCompletedCheckRuns($checks, $conclusion, [])`,
		},
		"$hasOnlyCompletedCheckRuns with three variables": {
			arg:     `$hasOnlyCompletedCheckRuns($checks, $conclusion, $ignoredConclusions)`,
			wantVal: `$hasOnlyCompletedCheckRuns($checks, $conclusion, $ignoredConclusions)`,
		},
		"join empty array with empty separator": {
			arg:     `$join([])`,
			wantVal: `$join([], " ")`,
		},
		"join empty array with non-empty separator": {
			arg:     `$join([], ", ")`,
			wantVal: `$join([], ", ")`,
		},
		"join non-empty array with empty separator": {
			arg:     `$join(["a", "b"])`,
			wantVal: `$join(["a", "b"], " ")`,
		},
		"join non-empty array with non-empty separator": {
			arg:     `$join(["a", "b"], ", ")`,
			wantVal: `$join(["a", "b"], ", ")`,
		},
		"join group with empty separator": {
			arg:     `$join($group("a"))`,
			wantVal: `$join($group("a"), " ")`,
		},
		"join group with non-empty separator": {
			arg:     `$join($group("a"), ", ")`,
			wantVal: `$join($group("a"), ", ")`,
		},
		"join team with empty separator": {
			arg:     `$join($team("a"))`,
			wantVal: `$join($team("a"), " ")`,
		},
		"join team with non-empty separator": {
			arg:     `$$join($team("a"), ", ")`,
			wantVal: `$$join($team("a"), ", ")`,
		},
		"join assignees with empty separator": {
			arg:     `$join($assignees())`,
			wantVal: `$join($assignees(), " ")`,
		},
		"join assignees with non-empty separator": {
			arg:     `$join($assignees(), ", ")`,
			wantVal: `$join($assignees(), ", ")`,
		},
		"join in sprintf with no separator": {
			arg:     `$sprintf("hello: %s", [$join(["test", "test2"])])`,
			wantVal: `$sprintf("hello: %s", [$join(["test", "test2"], " ")])`,
		},
		"join in sprintf with separator": {
			arg:     `$sprintf("hello: %s", [$join(["test", "test2"], " - ")])`,
			wantVal: `$sprintf("hello: %s", [$join(["test", "test2"], " - ")])`,
		},
		"join with one variable": {
			arg:     `$join($joinArray)`,
			wantVal: `$join($joinArray, " ")`,
		},
		"join with two variables": {
			arg:     `$join($joinArray, $joinSeparator)`,
			wantVal: `$join($joinArray, $joinSeparator)`,
		},
		"join with first variable": {
			arg:     `$join($joinArray, ", ")`,
			wantVal: `$join($joinArray, ", ")`,
		},
		"join with last variable": {
			arg:     `$join($assignees(), $joinSeparator)`,
			wantVal: `$join($assignees(), $joinSeparator)`,
		},
		"approve": {
			arg:     `$approve()`,
			wantVal: `$approve("")`,
		},
		"approve with comment": {
			arg:     `$approve("test")`,
			wantVal: `$approve("test")`,
		},
		"approve with empty comment": {
			arg:     `$approve("")`,
			wantVal: `$approve("")`,
		},
		"approve with variable": {
			arg:     `$approve($approveComment)`,
			wantVal: `$approve($approveComment)`,
		},
		"assign code author reviewer empty args": {
			arg:     `$assignCodeAuthorReviewers()`,
			wantVal: `$assignCodeAuthorReviewers(1, [], 0)`,
		},
		"assign code author reviewer total provided": {
			arg:     `$assignCodeAuthorReviewers(5)`,
			wantVal: `$assignCodeAuthorReviewers(5, [], 0)`,
		},
		"assign code author reviewer total and excluded provided": {
			arg:     `$assignCodeAuthorReviewers(5, ["john", "jane"])`,
			wantVal: `$assignCodeAuthorReviewers(5, ["john", "jane"], 0)`,
		},
		"assign code author reviewer total, excluded and max reviews provided": {
			arg:     `$assignCodeAuthorReviewers(5, ["john", "jane"], 3)`,
			wantVal: `$assignCodeAuthorReviewers(5, ["john", "jane"], 3)`,
		},
		"assign code author reviewer total, excluded as group and max reviews provided": {
			arg:     `$assignCodeAuthorReviewers(1, $group("excluded_reviewers"), 2)`,
			wantVal: `$assignCodeAuthorReviewers(1, $group("excluded_reviewers"), 2)`,
		},
		"assign code author reviewer total and excluded provided as team": {
			arg:     `$assignCodeAuthorReviewers(1, $team("reviewers"))`,
			wantVal: `$assignCodeAuthorReviewers(1, $team("reviewers"), 0)`,
		},
		"assign code author reviewer with one variable": {
			arg:     `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal)`,
			wantVal: `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, [], 0)`,
		},
		"assign code author reviewer with two variables": {
			arg:     `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, $excludedReviewers)`,
			wantVal: `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, $excludedReviewers, 0)`,
		},
		"assign code author reviewer with three variables": {
			arg:     `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, $excludedReviewers, $maxReviews)`,
			wantVal: `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, $excludedReviewers, $maxReviews)`,
		},
		"assign code author reviewer with first variable": {
			arg:     `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, [], 0)`,
			wantVal: `$assignCodeAuthorReviewers($assignCodeAuthorReviewersTotal, [], 0)`,
		},
		"assign code author reviewer with middle variable": {
			arg:     `$assignCodeAuthorReviewers(1, $excludedReviewers, 1)`,
			wantVal: `$assignCodeAuthorReviewers(1, $excludedReviewers, 1)`,
		},
		"assign code author reviewer with last variable": {
			arg:     `$assignCodeAuthorReviewers(1, ["john", "jane"], $maxReviews)`,
			wantVal: `$assignCodeAuthorReviewers(1, ["john", "jane"], $maxReviews)`,
		},
		"addReviewersBasedOnCodeAuthor empty args": {
			arg:     `$addReviewersBasedOnCodeAuthor()`,
			wantVal: `$addReviewersBasedOnCodeAuthor(1, [], 0)`,
		},
		"addReviewersBasedOnCodeAuthor total provided": {
			arg:     `$addReviewersBasedOnCodeAuthor(5)`,
			wantVal: `$addReviewersBasedOnCodeAuthor(5, [], 0)`,
		},
		"addReviewersBasedOnCodeAuthor total and excluded provided": {
			arg:     `$addReviewersBasedOnCodeAuthor(5, ["john", "jane"])`,
			wantVal: `$addReviewersBasedOnCodeAuthor(5, ["john", "jane"], 0)`,
		},
		"addReviewersBasedOnCodeAuthor total, excluded and max reviews provided": {
			arg:     `$addReviewersBasedOnCodeAuthor(5, ["john", "jane"], 3)`,
			wantVal: `$addReviewersBasedOnCodeAuthor(5, ["john", "jane"], 3)`,
		},
		"addReviewersBasedOnCodeAuthor total, excluded as group and max reviews provided": {
			arg:     `$addReviewersBasedOnCodeAuthor(1, $group("excluded_reviewers"), 2)`,
			wantVal: `$addReviewersBasedOnCodeAuthor(1, $group("excluded_reviewers"), 2)`,
		},
		"addReviewersBasedOnCodeAuthor total and excluded provided as team": {
			arg:     `$addReviewersBasedOnCodeAuthor(1, $team("reviewers"))`,
			wantVal: `$addReviewersBasedOnCodeAuthor(1, $team("reviewers"), 0)`,
		},
		"addReviewersBasedOnCodeAuthor with one variable": {
			arg:     `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal)`,
			wantVal: `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, [], 0)`,
		},
		"addReviewersBasedOnCodeAuthor with two variables": {
			arg:     `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, $excludedReviewers)`,
			wantVal: `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, $excludedReviewers, 0)`,
		},
		"addReviewersBasedOnCodeAuthor with three variables": {
			arg:     `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, $excludedReviewers, $maxReviews)`,
			wantVal: `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, $excludedReviewers, $maxReviews)`,
		},
		"addReviewersBasedOnCodeAuthor with first variable": {
			arg:     `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, [], 0)`,
			wantVal: `$addReviewersBasedOnCodeAuthor($addReviewersBasedOnCodeAuthorTotal, [], 0)`,
		},
		"addReviewersBasedOnCodeAuthor with middle variable": {
			arg:     `$addReviewersBasedOnCodeAuthor(1, $excludedReviewers, 1)`,
			wantVal: `$addReviewersBasedOnCodeAuthor(1, $excludedReviewers, 1)`,
		},
		"addReviewersBasedOnCodeAuthor with last variable": {
			arg:     `$addReviewersBasedOnCodeAuthor(1, ["john", "jane"], $maxReviews)`,
			wantVal: `$addReviewersBasedOnCodeAuthor(1, ["john", "jane"], $maxReviews)`,
		},
		"has any check run completed with no args": {
			arg:     `$hasAnyCheckRunCompleted()`,
			wantVal: `$hasAnyCheckRunCompleted([], [])`,
		},
		"has any check run completed with one arg": {
			arg:     `$hasAnyCheckRunCompleted($group("build"))`,
			wantVal: `$hasAnyCheckRunCompleted($group("build"), [])`,
		},
		"has any check run completed with two args": {
			arg:     `$hasAnyCheckRunCompleted(["build"], ["skipped"])`,
			wantVal: `$hasAnyCheckRunCompleted(["build"], ["skipped"])`,
		},
		"has any check run completed with group args": {
			arg:     `$hasAnyCheckRunCompleted($group("ignored"), $group("conclusions"))`,
			wantVal: `$hasAnyCheckRunCompleted($group("ignored"), $group("conclusions"))`,
		},
		"has any check run completed with first variable": {
			arg:     `$hasAnyCheckRunCompleted($ignoredChecks, [])`,
			wantVal: `$hasAnyCheckRunCompleted($ignoredChecks, [])`,
		},
		"has any check run completed with last variable": {
			arg:     `$hasAnyCheckRunCompleted([], $conclusionChecks)`,
			wantVal: `$hasAnyCheckRunCompleted([], $conclusionChecks)`,
		},
		"has any check run completed with all variables": {
			arg:     `$hasAnyCheckRunCompleted($ignoredChecks, $conclusionChecks)`,
			wantVal: `$hasAnyCheckRunCompleted($ignoredChecks, $conclusionChecks)`,
		},
		"has any check run completed with one nested function call and one variable": {
			arg:     `$hasAnyCheckRunCompleted($group("checks), $conclusionChecks)`,
			wantVal: `$hasAnyCheckRunCompleted($group("checks), $conclusionChecks)`,
		},
		"assign assignee with users list provided": {
			arg:     `$assignAssignees(["john", "mary"])`,
			wantVal: `$assignAssignees(["john", "mary"], 10)`,
		},
		"assign assignee with users list and total provided": {
			arg:     `$assignAssignees(["john", "mary"], 1)`,
			wantVal: `$assignAssignees(["john", "mary"], 1)`,
		},
		"assign assignee with group": {
			arg:     `$assignAssignees($group("test"))`,
			wantVal: `$assignAssignees($group("test"), 10)`,
		},
		"assign assignee with group and total provided": {
			arg:     `$assignAssignees($group("test"), 1)`,
			wantVal: `$assignAssignees($group("test"), 1)`,
		},
		"assign assignee with team": {
			arg:     `$assignAssignees($team("test"))`,
			wantVal: `$assignAssignees($team("test"), 10)`,
		},
		"assign assignee with team and total provided": {
			arg:     `$assignAssignees($team("test"), 1)`,
			wantVal: `$assignAssignees($team("test"), 1)`,
		},
		"assign assignee with one variable": {
			arg:     `$assignAssignees($assignAssigneesUsers, 5)`,
			wantVal: `$assignAssignees($assignAssigneesUsers, 5)`,
		},
		"assign asignees with two variables": {
			arg:     `$assignAssignees($assignAssigneesUsers, $assignAssigneesTotal)`,
			wantVal: `$assignAssignees($assignAssigneesUsers, $assignAssigneesTotal)`,
		},
		"assign asignees with last variable": {
			arg:     `$assignAssignees($team("test"), $assignAssigneesTotal)`,
			wantVal: `$assignAssignees($team("test"), $assignAssigneesTotal)`,
		},
		"assign asignees with first argument array and second variable": {
			arg:     `$assignAssignees(["john", "jane"], $assignAssigneesTotal)`,
			wantVal: `$assignAssignees(["john", "jane"], $assignAssigneesTotal)`,
		},
		"addAssignees with users list provided": {
			arg:     `$addAssignees(["john", "mary"])`,
			wantVal: `$addAssignees(["john", "mary"], 10)`,
		},
		"addAssignees with users list and total provided": {
			arg:     `$addAssignees(["john", "mary"], 1)`,
			wantVal: `$addAssignees(["john", "mary"], 1)`,
		},
		"addAssignees with group": {
			arg:     `$addAssignees($group("test"))`,
			wantVal: `$addAssignees($group("test"), 10)`,
		},
		"addAssignees with group and total provided": {
			arg:     `$addAssignees($group("test"), 1)`,
			wantVal: `$addAssignees($group("test"), 1)`,
		},
		"addAssignees with team": {
			arg:     `$addAssignees($team("test"))`,
			wantVal: `$addAssignees($team("test"), 10)`,
		},
		"addAssignees with team and total provided": {
			arg:     `$addAssignees($team("test"), 1)`,
			wantVal: `$addAssignees($team("test"), 1)`,
		},
		"addAssignees with one variable": {
			arg:     `$addAssignees($addAssigneesUsers, 5)`,
			wantVal: `$addAssignees($addAssigneesUsers, 5)`,
		},
		"addAssignees with two variables": {
			arg:     `$addAssignees($addAssigneesUsers, $addAssigneesTotal)`,
			wantVal: `$addAssignees($addAssigneesUsers, $addAssigneesTotal)`,
		},
		"addAssignees with last variable": {
			arg:     `$addAssignees($team("test"), $addAssigneesTotal)`,
			wantVal: `$addAssignees($team("test"), $addAssigneesTotal)`,
		},
		"addAssignees with first argument array and second variable": {
			arg:     `$addAssignees(["john", "jane"], $addAssigneesTotal)`,
			wantVal: `$addAssignees(["john", "jane"], $addAssigneesTotal)`,
		},
		"hasCodeWithoutSemanticChanges with no argument": {
			arg:     `$hasCodeWithoutSemanticChanges()`,
			wantVal: `$hasCodeWithoutSemanticChanges([])`,
		},
		"hasCodeWithoutSemanticChanges with default argument": {
			arg:     `$hasCodeWithoutSemanticChanges([])`,
			wantVal: `$hasCodeWithoutSemanticChanges([])`,
		},
		"hasCodeWithoutSemanticChanges with argument": {
			arg:     `$hasCodeWithoutSemanticChanges(["*.md", "*.txt"])`,
			wantVal: `$hasCodeWithoutSemanticChanges(["*.md", "*.txt"])`,
		},
		"hasCodeWithoutSemanticChanges with variable": {
			arg:     `$hasCodeWithoutSemanticChanges($hasCodeWithoutSemanticChangesPaths)`,
			wantVal: `$hasCodeWithoutSemanticChanges($hasCodeWithoutSemanticChangesPaths)`,
		},
		"containsOnlyCodeWithoutSemanticChanges with no argument": {
			arg:     `$containsOnlyCodeWithoutSemanticChanges()`,
			wantVal: `$containsOnlyCodeWithoutSemanticChanges([])`,
		},
		"containsOnlyCodeWithoutSemanticChanges with default argument": {
			arg:     `$containsOnlyCodeWithoutSemanticChanges([])`,
			wantVal: `$containsOnlyCodeWithoutSemanticChanges([])`,
		},
		"containsOnlyCodeWithoutSemanticChanges with argument": {
			arg:     `$containsOnlyCodeWithoutSemanticChanges(["*.md", "*.txt"])`,
			wantVal: `$containsOnlyCodeWithoutSemanticChanges(["*.md", "*.txt"])`,
		},
		"containsOnlyCodeWithoutSemanticChanges with variable": {
			arg:     `$containsOnlyCodeWithoutSemanticChanges($containsOnlyCodeWithoutSemanticChangesPaths)`,
			wantVal: `$containsOnlyCodeWithoutSemanticChanges($containsOnlyCodeWithoutSemanticChangesPaths)`,
		},
		"summarize alias": {
			arg:     `$summarize()`,
			wantVal: `$robinSummarize("default", "openai-gpt-4")`,
		},
		"assignReviewer with one variable": {
			arg:     `$assignReviewer($reviewers)`,
			wantVal: `$assignReviewer($reviewers, 99, "reviewpad")`,
		},
		"assignReviewer with two variables": {
			arg:     `$assignReviewer($reviewers, $maxReviewers)`,
			wantVal: `$assignReviewer($reviewers, $maxReviewers, "reviewpad")`,
		},
		"assignReviewer with three variables": {
			arg:     `$assignReviewer($reviewers, $maxReviewers, $strategy)`,
			wantVal: `$assignReviewer($reviewers, $maxReviewers, $strategy)`,
		},
		"assignReviewer with one literal": {
			arg:     `$assignReviewer(["john"])`,
			wantVal: `$assignReviewer(["john"], 99, "reviewpad")`,
		},
		"assignReviewer with two literals": {
			arg:     `$assignReviewer(["john"], 1)`,
			wantVal: `$assignReviewer(["john"], 1, "reviewpad")`,
		},
		"assignReviewer with three literals": {
			arg:     `$assignReviewer(["john"], 1, "reviewpad")`,
			wantVal: `$assignReviewer(["john"], 1, "reviewpad")`,
		},
		"assignReviewer with one literal and one variable": {
			arg:     `$assignReviewer(["john", "jane"], $maxReviewers)`,
			wantVal: `$assignReviewer(["john", "jane"], $maxReviewers, "reviewpad")`,
		},
		"assignReviewer with one literal and two variables": {
			arg:     `$assignReviewer(["john", "jane"], $maxReviewers, $strategy)`,
			wantVal: `$assignReviewer(["john", "jane"], $maxReviewers, $strategy)`,
		},
		"addReviewers with one variable": {
			arg:     `$addReviewers($reviewers)`,
			wantVal: `$addReviewers($reviewers, 99, "reviewpad")`,
		},
		"addReviewers with two variables": {
			arg:     `$addReviewers($reviewers, $maxReviewers)`,
			wantVal: `$addReviewers($reviewers, $maxReviewers, "reviewpad")`,
		},
		"addReviewers with three variables": {
			arg:     `$addReviewers($reviewers, $maxReviewers, $strategy)`,
			wantVal: `$addReviewers($reviewers, $maxReviewers, $strategy)`,
		},
		"addReviewers with one literal": {
			arg:     `$addReviewers(["john"])`,
			wantVal: `$addReviewers(["john"], 99, "reviewpad")`,
		},
		"addReviewers with two literals": {
			arg:     `$addReviewers(["john"], 1)`,
			wantVal: `$addReviewers(["john"], 1, "reviewpad")`,
		},
		"addReviewers with three literals": {
			arg:     `$addReviewers(["john"], 1, "reviewpad")`,
			wantVal: `$addReviewers(["john"], 1, "reviewpad")`,
		},
		"addReviewers with one literal and one variable": {
			arg:     `$addReviewers(["john", "jane"], $maxReviewers)`,
			wantVal: `$addReviewers(["john", "jane"], $maxReviewers, "reviewpad")`,
		},
		"addReviewers with one literal and two variables": {
			arg:     `$addReviewers(["john", "jane"], $maxReviewers, $strategy)`,
			wantVal: `$addReviewers(["john", "jane"], $maxReviewers, $strategy)`,
		},
		"getReviewers": {
			arg:     `$getReviewers()`,
			wantVal: `$getReviewers("")`,
		},
		"getReviewers approved": {
			arg:     `$getReviewers("APPROVED")`,
			wantVal: `$getReviewers("APPROVED")`,
		},
		// TODO: test addDefaultTotalRequestedReviewers
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotVal := transformAladinoExpression(test.arg)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
