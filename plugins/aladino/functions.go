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
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
	"github.com/shurcooL/githubv4"
)

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

func base() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: baseCode,
	}
}

func baseCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetBase().GetRef()), nil
}

func commitCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: commitCountCode,
	}
}

func commitCountCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildIntValue(*e.GetPullRequest().Commits), nil
}

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

func description() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: descriptionCode,
	}
}

func descriptionCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetBody()), nil
}

func fileCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: fileCountCode,
	}
}

func fileCountCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	patch := e.GetPatch()
	return aladino.BuildIntValue(len(patch)), nil
}

func hasCodePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasCodePatternCode,
	}
}

func patchHasCodePatternCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	arg := args[0].(*aladino.StringValue)
	patch := e.GetPatch()

	for _, file := range patch {
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
	for fp := range patch {
		fpExt := utils.FileExt(fp)
		normalizedExt := strings.ToLower(fpExt)

		if _, ok := extensionSet[normalizedExt]; !ok {
			return aladino.BuildFalseValue(), nil
		}
	}

	return aladino.BuildTrueValue(), nil
}

func hasFileName() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasFileNameCode,
	}
}

func patchHasFileNameCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	fileNameStr := args[0].(*aladino.StringValue)

	patch := e.GetPatch()
	for fp := range patch {
		if fp == fileNameStr.Val {
			return aladino.BuildTrueValue(), nil
		}
	}

	return aladino.BuildFalseValue(), nil
}

func hasFilePattern() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasFilePatternCode,
	}
}

func patchHasFilePatternCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	filePatternRegex := args[0].(*aladino.StringValue)

	patch := e.GetPatch()
	for fp := range patch {
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
		if len(ghCommit.Parents) > 1 {
			return aladino.BuildBoolValue(false), nil
		}
	}

	return aladino.BuildBoolValue(true), nil
}

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

func head() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: headCode,
	}
}

func headCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetHead().GetRef()), nil
}

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

func title() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code: titleCode,
	}
}

func titleCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	return aladino.BuildStringValue(e.GetPullRequest().GetTitle()), nil
}

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

func group() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code: groupCode,
	}
}

func groupCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	groupName := args[0].(*aladino.StringValue).Val

	if val, ok := e.GetRegisterMap()[groupName]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("getGroup: no group with name %v in state %+q", groupName, e.GetRegisterMap())
}

func rule() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: ruleCode,
	}
}

func ruleCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	ruleName := args[0].(*aladino.StringValue).Val

	internalRuleName := aladino.BuildInternalRuleName(ruleName)

	if spec, ok := e.GetRegisterMap()[internalRuleName]; ok {
		specRaw := spec.(*aladino.StringValue).Val
		result, err := aladino.EvalExpr(e, "patch", specRaw)
		if err != nil {
			return nil, err
		}
		return aladino.BuildBoolValue(result), nil
	}

	return nil, fmt.Errorf("$rule: no rule with name %v in state %+q", ruleName, e.GetRegisterMap())
}

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
