// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"fmt"
	"strings"
	"time"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/lang/aladino"
	"github.com/reviewpad/reviewpad/utils"
	"github.com/shurcooL/githubv4"
)

/*
reviewpad-an: builtin-docs

# Functions
______________

Reviewpad functions allow to query data from a `pull request` or `organization` in order to act on it.

The functions are organized into 4 categories:
- **[Pull Request](#pull-request)** - Functions to query pull request data.
- **[Organization](#organization)** - Functions to query organization data.
- **[User](#user)** - Functions to query user's data.
- **[Utilities](#utilities)** - Functions to help act on the queried data.
- **[Engine](#engine)** - Functions used to work with `reviewpad.yml` file.

*/

/*
reviewpad-an: builtin-docs

## Pull Request
______________

Set of functions to get pull request details.
*/

/*
reviewpad-an: builtin-docs

## assignees
______________

**Description**:

Returns the list of GitHub user login that are assigned to the pull request.

**Parameters**:

*none*

**Return value**:

`[]string`

The list of GitHub user login that are assigned to the pull request.

**Examples**:

```yml
$assignees()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    assignedToTechLead:
        kind: patch
        description: Verifies if pull request was assigned only to a specific tech lead
        spec: $assignees() == ["john"]
```
*/
func assignees() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: assigneesCode,
	}
}

func assigneesCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	ghAssignees := e.GetPullRequest().Assignees
	assignees := make([]aladino.Value, len(ghAssignees))

	for i, ghAssignee := range ghAssignees {
		assignees[i] = aladino.BuildStringValue(ghAssignee.GetLogin())
	}

	return aladino.BuildArrayValue(assignees), nil
}

/*
reviewpad-an: builtin-docs

## author
______________

**Description**:

Retrieves the pull request author GitHub login.

**Parameters**:

*none*

**Return value**:

`string`

The GitHub login of the pull request author.

**Examples**:

```yml
$author()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    isAuthoredByTechLead:
        kind: patch
        description: Verifies if authored by tech lead
        spec: $author() == "john"
```
*/
func author() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: authorCode,
	}
}

func authorCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	authorLogin := e.GetPullRequest().User.GetLogin()
	return aladino.BuildStringValue(authorLogin), nil
}

/*
reviewpad-an: builtin-docs

## base
______________

**Description**:

Returns the name of the branch the pull request should be pulled into.

**Parameters**:

*none*

**Return value**:

`string`

The name of the branch the pull request should be pulled into.

**Examples**:

```yml
$base()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    shouldNotifyTechLead:
        kind: patch
        description: Verifies if pull request is going to be pulled into "features" branch
        spec: $base() == "features"
```
*/
func base() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: baseCode,
	}
}

func baseCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetBase().GetRef()), nil
}

/*
reviewpad-an: builtin-docs

## commitCount
______________

**Description**:

Returns the total number of commits made into the pull request.

**Parameters**:

*none*

**Return value**:

`int`

The total number of commits in the pull request.

**Examples**:

```yml
$commitCount()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasTooManyCommits:
        kind: patch
        description: Verifies if it has than 3 commits
        spec: $commitCount() > 3
```
*/
func commitCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: commitCountCode,
	}
}

func commitCountCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildIntValue(*e.GetPullRequest().Commits), nil
}

/*
reviewpad-an: builtin-docs

## commits
______________

**Description**:

Returns the list of commit messages of the pull request.

**Parameters**:

*none*

**Return value**:

`[]string`

The list of commit messages of the pull request.

**Examples**:

```yml
$commits()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    skipCIOnCommitMention:
        kind: patch
        description: Verifies if any of the commit messages of the pull request contains "skip-ci"
        spec: $elemContains("skip-ci", $commits())
```
*/
func commits() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: commitsCode,
	}
}

func commitsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	ghCommits, _, err := e.GetClient().PullRequests.ListCommits(e.GetCtx(), owner, repo, prNum, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	commitMessages := make([]aladino.Value, len(ghCommits))
	for i, ghCommit := range ghCommits {
		commitMessages[i] = aladino.BuildStringValue(ghCommit.Commit.GetMessage())
	}

	return aladino.BuildArrayValue(commitMessages), nil
}

/*
reviewpad-an: builtin-docs

## createdAt
______________

**Description**:

Returns the time the pull request was created at.

**Parameters**:

*none*

**Return value**:

`int64`

The number of seconds elapsed since January 1, 1970 UTC.

**Examples**:

```yml
$createdAt()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    wasCreatedOnApril:
        kind: patch
        description: Verifies if the pull request was created on the April 14th of 2011 at 16:00:49
        spec: $createdAt() == "2011-04-14T16:00:49Z
```
*/
func createdAt() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: createdAtCode,
	}
}

func createdAtCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	createdAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", e.GetPullRequest().GetCreatedAt().String())
	if err != nil {
		return nil, err
	}

	return aladino.BuildIntValue(int(createdAtTime.Unix())), nil
}

/*
reviewpad-an: builtin-docs

## description
______________

**Description**:

Returns the description of the pull request.

**Parameters**:

*none*

**Return value**:

`string`

The description of the pull request.

**Examples**:

```yml
$description()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasDescription:
        kind: patch
        description: Verifies if the pull request description is "Testing description"
        spec: $description() == "Testing description"
```
*/
func description() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: descriptionCode,
	}
}

func descriptionCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetBody()), nil
}

/*
reviewpad-an: builtin-docs

## fileCount
______________

**Description**:

Returns the total number of files changed in the patch.

**Parameters**:

*none*

**Return value**:

`int`

The total number of files changed in the patch.

**Examples**:

```yml
$filesCount()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    changesTooManyFiles:
        kind: patch
        description: Verifies if it has than 3 files
        spec: $fileCount() > 3
```
*/
func fileCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: fileCountCode,
	}
}

func fileCountCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	patch := e.GetPatch()
	return aladino.BuildIntValue(len(*patch)), nil
}

/*
reviewpad-an: builtin-docs

## hasCodePattern
______________

**Description**:

Verifies if the patch matches the provided code pattern, returning `true` or `false` as appropriate.

The code pattern needs to be a compilable regular expression.

**Parameters**:

| variable       | type   | description                        |
| -------------- | ------ | ---------------------------------- |
| `queryPattern` | string | query pattern to look for on patch |

**Return value**:

`boolean`

Returns `true` if the patch matches the code query, `false` otherwise.

**Examples**:

```yml
$codePattern("placeBet\(.*\)")
```

A `reviewpad.yml.yml` example:

```yml
rules:
    usesPlaceBet:
        kind: patch
        description: Verifies if uses placeBet
        spec: $codePattern("placeBet\(.*\)")
```
*/
func hasCodePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasCodePatternCode,
	}
}

func patchHasCodePatternCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	arg := args[0].(*aladino.StringValue)
	patch := e.GetPatch()

	for _, file := range *patch {
		if file == nil {
			continue
		}

		isMatch, err := file.Query(arg.Val)
		if err != nil {
			return nil, err
		}

		if isMatch {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}

/*
reviewpad-an: builtin-docs

## hasFileExtensions
______________

**Description**:

Determines whether all the extensions of the changed files on the patch are included on the provided list of file extensions, returning `true` or `false` as appropriate.

Each extension provided on the list needs to be a [glob](https://en.wikipedia.org/wiki/Glob_(programming)).

**Parameters**:

| variable     | type     | description                 |
| ------------ | -------- | --------------------------- |
| `extensions` | []string | list of all file extensions |

**Return value**:

`boolean`

Returns `true` if all file extensions in the patch are included in the list, `false` otherwise.

**Examples**:

```yml
$hasFileExtensions([".test.ts"])
```

A `reviewpad.yml.yml` example:

```yml
rules:
    changesAreOnlyTests:
        kind: patch
        description: Verifies if changes are only on test files
        spec: $hasFileExtensions([".test.ts"])
```
*/
func hasFileExtensions() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildBoolType()),
		Code: patchHasFileExtensionsCode,
	}
}

func patchHasFileExtensionsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	argExtensions := args[0].(*aladino.ArrayValue)

	extensionSet := make(map[string]bool, len(argExtensions.Vals))
	for _, argExt := range argExtensions.Vals {
		argStringVal := argExt.(*aladino.StringValue)

		normalizedStr := strings.ToLower(argStringVal.Val)
		extensionSet[normalizedStr] = true
	}

	patch := e.GetPatch()
	for fp := range *patch {
		fpExt := utils.FileExt(fp)
		normalizedExt := strings.ToLower(fpExt)

		if _, ok := extensionSet[normalizedExt]; !ok {
			return aladino.BuildFalseValue(), nil
		}
	}

	return aladino.BuildTrueValue(), nil
}

/*
reviewpad-an: builtin-docs

## hasFileName
______________

**Description**:

Determines whether the provided filename is among the files on patch, returning `true` or `false` as appropriate.

**Parameters**:

| variable   | type   | description                                        |
| ---------- | ------ | -------------------------------------------------- |
| `filename` | string | filename to look for in the patch. case sensitive. |

**Return value**:

`boolean`

Returns `true` if the patch has a file with the provided filename, `false` otherwise.

The provided filename and the filename on the patch need to be exactly the same in order to get a positive result.

**Examples**:

```yml
$hasFileName("placeBet.js")
```

A `reviewpad.yml.yml` example:

```yml
rules:
    changesPlaceBet:
        kind: patch
        description: Verifies if changes place bet file
        spec: $hasFileName("placeBet.js")
```
*/
func hasFileName() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasFileNameCode,
	}
}

func patchHasFileNameCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	fileNameStr := args[0].(*aladino.StringValue)

	patch := e.GetPatch()
	for fp := range *patch {
		if fp == fileNameStr.Val {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}

/*
reviewpad-an: builtin-docs

## hasFilePattern
______________

**Description**:

Determines whether the provided file pattern matches any of the files in the patch, returning `true` or `false` as appropriate.

The file pattern needs to be a [glob](https://en.wikipedia.org/wiki/Glob_(programming)).

**Parameters**:

| variable      | type   | description                            |
| ------------- | ------ | -------------------------------------- |
| `filePattern` | string | file pattern glob to look for on patch |

**Return value**:

`boolean`

Returns `true` if any of the files on patch matches the provided file pattern, `false` otherwise.

**Examples**:

```yml
$hasFilePattern("src/transactions/**")
```

A `reviewpad.yml.yml` example:

```yml
rules:
    changesTransactions:
        kind: patch
        description: Verifies if changes transactions
        spec: $hasFilePattern("src/transactions/**")
```
*/
func hasFilePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasFilePatternCode,
	}
}

func patchHasFilePatternCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	filePatternRegex := args[0].(*aladino.StringValue)

	patch := e.GetPatch()
	for fp := range *patch {
		re, err := doublestar.Match(filePatternRegex.Val, fp)
		if err != nil {
			return aladino.BuildFalseValue(), err
		}
		if re {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}

/*
reviewpad-an: builtin-docs

## hasLinearHistory
______________

**Description**:

Checks if a pull request has a linear history.

A linear history is simply a Git history in which all commits come after one another, i.e., you will not find any merges of branches with independent commit histories.

**Parameters**:

*none*

**Return value**:

`boolean`

Returns `true` if it has a linear history. `false` otherwise.

**Examples**:

```yml
$hasLinearHistory()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasLinearHistory:
        kind: patch
        description: Verifies if the pull request has a linear history
        spec: $hasLinearHistory()
```
*/
func hasLinearHistory() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasLinearHistoryCode,
	}
}

func hasLinearHistoryCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	ghCommits, _, err := e.GetClient().PullRequests.ListCommits(e.GetCtx(), owner, repo, prNum, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ghCommit := range ghCommits {
		if !utils.HasLinearHistory(ghCommit.Commit) {
			return aladino.BuildBoolValue(false), nil
		}
	}

	return aladino.BuildBoolValue(true), nil
}

/*
reviewpad-an: builtin-docs

## hasLinkedIssues
______________

**Description**:

Checks if a pull request has associated issues that might be closed by it.

**Parameters**:

*none*

**Return value**:

`boolean`

Returns `true` if it has linked issues. `false` otherwise.

**Examples**:

```yml
$hasLinkedIssues()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasLinkedIssues:
        kind: patch
        description: Verifies if the pull request has linked issues
        spec: $hasLinkedIssues()
```
*/
func hasLinkedIssues() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: hasLinkedIssuesCode,
	}
}

func hasLinkedIssuesCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	var pullRequestQuery struct {
		Repository struct {
			PullRequest struct {
				ClosingIssuesReferences struct {
					TotalCount githubv4.Int
				}
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	varGQLPullRequestQuery := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(prNum),
	}

	err := e.GetClientGQL().Query(e.GetCtx(), &pullRequestQuery, varGQLPullRequestQuery)

	if err != nil {
		return nil, err
	}

	return aladino.BuildBoolValue(pullRequestQuery.Repository.PullRequest.ClosingIssuesReferences.TotalCount > 0), nil
}

/*
reviewpad-an: builtin-docs

## head
______________

**Description**:

Returns the name of the branch where the pull request changes are implemented.

**Parameters**:

*none*

**Return value**:

`string`

The name of the branch where the pull request changes are implemented.

**Examples**:

```yml
$head()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    changesImplementedInDevelopmentBranch:
        kind: patch
        description: Verifies if pull request changes are implemented in the "development" branch
        spec: $head() == "development"
```
*/
func head() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: headCode,
	}
}

func headCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetHead().GetRef()), nil
}

/*
reviewpad-an: builtin-docs

## isDraft
______________

**Description**:

Verifies whether the pull request is `Draft`, returning `true` or `false` as appropriate.

To know more about [GitHub Draft pull request](https://github.blog/2019-02-14-introducing-draft-pull-requests/).

**Parameters**:

*none*

**Return value**:

`boolean`

A boolean al.Value which is `true` if the pull request is `Draft`, `false` otherwise.

**Examples**:

```yml
$isDraft()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    isDraft:
        kind: patch
        description: Verifies if is Draft
        spec: $isDraft()
```
*/
func isDraft() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code: isDraftCode,
	}
}

func isDraftCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	if pullRequest == nil {
		return nil, fmt.Errorf("isDraft: pull request is nil")
	}

	return aladino.BuildBoolValue(pullRequest.GetDraft()), nil
}

/*
reviewpad-an: builtin-docs

## labels
______________

**Description**:

Returns the list of labels of the pull request.

**Parameters**:

*none*

**Return value**:

`[]string`

The list of labels of the pull request.

**Examples**:

```yml
$labels()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    onlyHasTestLabelAssigned:
        kind: patch
        description: Verifies if the pull request only has "test" label assigned
        spec: $labels() == ["test"]
```
*/
func labels() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: labelsCode,
	}
}

func labelsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	ghLabels := e.GetPullRequest().Labels
	labels := make([]aladino.Value, len(ghLabels))

	for i, ghLabel := range ghLabels {
		labels[i] = aladino.BuildStringValue(ghLabel.GetName())
	}

	return aladino.BuildArrayValue(labels), nil
}

/*
reviewpad-an: builtin-docs

## milestone
______________

**Description**:

Returns the milestone title associated to the pull request.

**Parameters**:

*none*

**Return value**:

`string`

The milestone title associated to the pull request.

**Examples**:

```yml
$milestone()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    isPartOfBugFixesMilestone:
        kind: patch
        description: Verifies if the pull request is associated with the bug fixes milestone
        spec: $milestone() == "Bug fixes"
```
*/
func milestone() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: milestoneCode,
	}
}

func milestoneCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	milestoneTitle := e.GetPullRequest().GetMilestone().GetTitle()
	return aladino.BuildStringValue(milestoneTitle), nil
}

/*
reviewpad-an: builtin-docs

## reviewers
______________

**Description**:

Returns the list of GitHub user login or team slug that were requested to review the pull request.

**Parameters**:

*none*

**Return value**:

`[]string`

The list of GitHub user login or team slug that were requested to review the pull request.

**Examples**:

```yml
$reviewers()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasRequestedReviewers:
        kind: patch
        description: Verifies if the pull request has reviewers
        spec: $reviewers() != []
```
*/
func reviewers() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: reviewersCode,
	}
}

func reviewersCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	usersReviewers := e.GetPullRequest().RequestedReviewers
	teamReviewers := e.GetPullRequest().RequestedTeams
	totalReviewers := len(usersReviewers) + len(teamReviewers)
	reviewersLogin := make([]aladino.Value, totalReviewers)

	for i, userReviewer := range usersReviewers {
		reviewersLogin[i] = aladino.BuildStringValue(userReviewer.GetLogin())
	}

	for i, teamReviewer := range teamReviewers {
		reviewersLogin[i+len(usersReviewers)] = aladino.BuildStringValue(teamReviewer.GetSlug())
	}

	return aladino.BuildArrayValue(reviewersLogin), nil
}

/*
reviewpad-an: builtin-docs

## size
______________

**Description**:

Returns the total amount of changed lines in the patch.

Any added or removed line is considered a change. For instance, the following patch will have a `size` of 2 (one line removed and one line added.)

```diff
function helloWorld() {
-   return "Hello"
+   return "Hello World"
}
```

**Parameters**:

*none*

**Return value**:

`int`

The sum of all changed lines in the patch.

**Examples**:

```yml
$size()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    isBigChange:
        kind: patch
        description: Verifies if change is big
        spec: $size() > 100
```
*/
func size() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: sizeCode,
	}
}

func sizeCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	size := e.GetPullRequest().GetAdditions() + e.GetPullRequest().GetDeletions()
	return aladino.BuildIntValue(size), nil
}

/*
reviewpad-an: builtin-docs

## title
______________

**Description**:

Returns the title of the pull request.

**Parameters**:

*none*

**Return value**:

`string`

The title of the pull request.

**Examples**:

```yml
$title()
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasTitle:
        kind: patch
        description: Verifies if the pull request title is "Test custom builtins"
        spec: $title() == "Test custom builtins"
```
*/
func title() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: titleCode,
	}
}

func titleCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetTitle()), nil
}

/*
reviewpad-an: builtin-docs

## Organization
______________

Set of functions to get organization details.
*/

/*
reviewpad-an: builtin-docs

## organization
______________

**Description**:

Lists all the members of the organization that owns the pull request.

If the authenticated user is an owner of the organization, this will return both concealed and public members, otherwise it will only return public members.

**Parameters**:

*none*

**Return value**:

`[]string`

The list of all members of the organization to where the pull request running against.

**Examples**:

```yml
$organization()
````

A `reviewpad.yml.yml` example:

```yml
rules:
    isAuthorFromOrganization:
        kind: patch
        description: Verifies if author belongs to organization
        spec: $isElementOf($author(), $organization())
```
*/
func organization() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: organizationCode,
	}
}

func organizationCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	orgName := e.GetPullRequest().Head.Repo.Owner.Login
	users, _, err := e.GetClient().Organizations.ListMembers(e.GetCtx(), *orgName, nil)
	if err != nil {
		return nil, err
	}
	elems := make([]aladino.Value, len(users))
	for i, user := range users {
		elems[i] = aladino.BuildStringValue(*user.Login)
	}

	return aladino.BuildArrayValue(elems), nil
}

/*
reviewpad-an: builtin-docs

## team
______________

**Description**:

Returns the members of a team and child teams.

To list members in a team, the team must be visible to the authenticated user.

| :warning: Requires a GitHub token :warning: |
|---------------------------------------------|

By default a GitHub action does not have permission to access organization members.

Because of that, in order for the function `team` to work we need to provide a GitHub token to the Reviewpad action.

[Please follow this link to know more](https://docs.reviewpad.com/docs/install-github-action-with-github-token).

**Parameters**:

| variable       | type   | description                          |
| -------------- | ------ | ------------------------------------ |
| `teamSlug`     | string | The slug of the team name on GitHub. |

**Return value**:

`[]string`

Returns the list of all team and child teams members GitHub login.

**Examples**:

```yml
$team("devops")
```

A `reviewpad.yml.yml` example:

```yml
rules:
  isAuthorByDevops:
    description: Verifies if author belongs to devops team
    kind: patch
    spec: $isElementOf($author(), $team("devops"))
```
*/
func team() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: teamCode,
	}
}

func teamCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	teamSlug := args[0].(*aladino.StringValue).Val
	orgName := e.GetPullRequest().Head.Repo.Owner.Login

	members, _, err := e.GetClient().Teams.ListTeamMembersBySlug(e.GetCtx(), *orgName, teamSlug, &github.TeamListTeamMembersOptions{})
	if err != nil {
		return nil, err
	}

	membersLogin := make([]aladino.Value, len(members))
	for i, member := range members {
		membersLogin[i] = aladino.BuildStringValue(*member.Login)
	}

	return aladino.BuildArrayValue(membersLogin), nil
}

/*
reviewpad-an: builtin-docs

## User

Set of functions to get user details.
*/

/*
reviewpad-an: builtin-docs

## totalCreatedPullRequests
______________

**Description**:

Returns the total number of pull requests created by the provided GitHub user login.

**Parameters**:

| variable    | type   | description           |
| ----------- | ------ | --------------------- |
| `userLogin` | string | the GitHub user login |

**Return value**:

`int`

The total number of created pull requests created by GitHub user login.

**Examples**:

```yml
$totalCreatedPullRequests($author())
```

A `reviewpad.yml.yml` example:

```yml
rules:
    isJunior:
        kind: patch
        description: Verifies if author is junior
        spec: $totalCreatedPullRequests($author()) < 3
```
*/
func totalCreatedPullRequests() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildIntType()),
		Code: totalCreatedPullRequestsCode,
	}
}

func totalCreatedPullRequestsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	devName := args[0].(*aladino.StringValue).Val

	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	issues, _, err := e.GetClient().Issues.ListByRepo(e.GetCtx(), owner, repo, &github.IssueListByRepoOptions{
		Creator: devName,
		State:   "all",
	})
	if err != nil {
		return nil, err
	}

	count := 0
	for _, issue := range issues {
		if issue.IsPullRequest() {
			count++
		}
	}

	return aladino.BuildIntValue(count), nil
}

/*
reviewpad-an: builtin-docs

## Utilities
______________

Set of functions to help handle the queried data.
*/

/*
reviewpad-an: builtin-docs

## append
______________

**Description**:

Appends elements to the end of a slice and returns the updated slice.

**Parameters**:

| variable       | type     | description                                  |
| -------------- | -------- | -------------------------------------------- |
| `slice`        | []string | slice that will have elements appended to it |
| `elements`     | []string | elements to be added to the end of the slice |

**Return value**:

`[]string`

Returns a new slice by appending the slices passed to it.

**Examples**:

```yml
$append(["a", "b"], ["c"])    # ["a", "b", "c"]
```

A `reviewpad.yml.yml` example:

```yml
groups:
    frontendAndBackendDevs:
        description: Frontend and backend developers
        kind: developers
        spec: $append($team("frontend"), $team("backend"))
rules:
    authoredByWebDeveloper:
        kind: patch
        description: Authored by web developers
        spec: $isElementOf($author(), $group("frontendAndBackendDevs"))
```
*/
func appendString() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: appendStringCode,
	}
}

func appendStringCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	slice1 := args[0].(*aladino.ArrayValue).Vals
	slice2 := args[1].(*aladino.ArrayValue).Vals

	return aladino.BuildArrayValue(append(slice1, slice2...)), nil
}

/*
reviewpad-an: builtin-docs

## contains
______________

**Description**:

Determines whether a text includes a certain sentence, returning `true` or `false` as appropriate.

**Parameters**:

| variable            | type          | description                 |
| ------------------- | ------------- | --------------------------- |
| `text`              | string        | The text to search in       |
| `searchSentence`    | string        | The sentence to search for  |

**Return value**:

`boolean`

Returns `true` if the al.Value searchSentence is found within the text, `false` otherwise.

**Examples**:

```yml
$contains("Testing string contains", "string contains")     #true
$contains("Testing string contains", "test")                #false
```

A `reviewpad.yml.yml` example:

```yml
rules:
    hasCustomKeywordInTitle:
        kind: patch
        description: Verifies if the pull request title has "custom" keyword
        spec: $contains($title(), "custom")
```
*/
func contains() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: containsCode,
	}
}

func containsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	str := args[0].(*aladino.StringValue).Val
	subString := args[1].(*aladino.StringValue).Val

	return aladino.BuildBoolValue(strings.Contains(str, subString)), nil
}

/*
reviewpad-an: builtin-docs

## isElementOf
______________

**Description**:

Determines whether a list includes a certain al.Value among its entries, returning `true` or `false` as appropriate.

**Parameters**:

| variable        | type      | description             |
| --------------- | --------- | ----------------------- |
| `searchElement` | literal   | The al.Value to search for |
| `list`          | []literal | The list to search in   |

**Return value**:

`boolean`

Returns `true` if the al.Value searchElement is found within the list, `false` otherwise.

**Examples**:

```yml
$isElementOf("john", ["maria", "john"])  # true
$isElementOf(3, [1, 2])                  # false
```

A `reviewpad.yml.yml` example:

```yml
rules:
  authoredByJunior:
    description: Verifies if author is junior
    kind: patch
    spec: $isElementOf($author(), $group("junior"))
```
*/
func isElementOf() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildBoolType()),
		Code: isElementOfCode,
	}
}

func isElementOfCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	member := args[0].(*aladino.StringValue)
	group := args[1].(*aladino.ArrayValue).Vals

	for _, groupMember := range group {
		if member.Equals(groupMember) {
			return aladino.BuildBoolValue(true), nil
		}
	}

	return aladino.BuildBoolValue(false), nil
}

/*
reviewpad-an: builtin-docs

## Engine
______________

Set of functions used to handle `reviewpad.yml` file.

This functions should be used to access and handle data declared into `reviewpad.yml`, e.g. `$group` to get a defined group.
*/

/*
reviewpad-an: builtin-docs

## group
______________

**Description**:

Lists all members that belong to the provided group. This group needs to be defined in the same `reviewpad.yml.yml` file.

`group` is a way to refer to a defined set of users in a short way.

**Parameters**:

| variable    | type   | description                            |
| ----------- | ------ | -------------------------------------- |
| `groupName` | string | the group name to list the member from |

**Return value**:

`[]string`

Returns all members from the group.

**Examples**:

```yml
$group("techLeads")
```

A `reviewpad.yml.yml` example:

```yml
groups:
  techLeads:
    description: Group with all tech leads
    kind: developers
    spec: '["john", "maria", "arthur"]'

rules:
  isAuthorByTechLead:
    description: Verifies if author is a tech lead
    kind: patch
    spec: $isElementOf($author(), $group("techLeads"))
```
*/
func group() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: groupCode,
	}
}

func groupCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	groupName := args[0].(*aladino.StringValue).Val

	if val, ok := (*e.GetRegisterMap())[groupName]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("getGroup: no group with name %v in state %+q", groupName, e.GetRegisterMap())
}

/*
 * Internal
 */

func filter() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{
				aladino.BuildArrayOfType(aladino.BuildStringType()),
				aladino.BuildFunctionType(
					[]aladino.Type{aladino.BuildStringType()},
					aladino.BuildBoolType(),
				),
			},
			aladino.BuildArrayOfType(aladino.BuildStringType()),
		),
		Code: filterCode,
	}
}

func filterCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	result := make([]aladino.Value, 0)
	elems := args[0].(*aladino.ArrayValue).Vals
	fn := args[1].(*aladino.FunctionValue).Fn

	for _, elem := range elems {
		fnResult := fn([]aladino.Value{elem}).(*aladino.BoolValue).Val
		if fnResult {
			result = append(result, elem)
		}
	}

	return aladino.BuildArrayValue(result), nil
}
