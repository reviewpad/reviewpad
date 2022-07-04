// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	actions "github.com/reviewpad/reviewpad/v2/plugins/aladino/actions"
	pullRequest_functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions/pullRequest"
	organization_functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions/organization"
	utils_functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions/utils"
	engine_functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions/engine"
	inner_functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions/inner"
	user_functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions/user"
)

// The documentation for the builtins is in:
// https://github.com/reviewpad/docs/blob/main/aladino/builtins.md
// This means that changes to the builtins need to be propagated to that document.

func PluginBuiltIns() *aladino.BuiltIns {
	return &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			// Pull Request
			"assignees":         pullRequest_functions.Assignees(),
			"author":            pullRequest_functions.Author(),
			"base":              pullRequest_functions.Base(),
			"commitCount":       pullRequest_functions.CommitCount(),
			"commits":           pullRequest_functions.Commits(),
			"createdAt":         pullRequest_functions.CreatedAt(),
			"description":       pullRequest_functions.Description(),
			"fileCount":         pullRequest_functions.FileCount(),
			"hasCodePattern":    pullRequest_functions.HasCodePattern(),
			"hasFileExtensions": pullRequest_functions.HasFileExtensions(),
			"hasFileName":       pullRequest_functions.HasFileName(),
			"hasFilePattern":    pullRequest_functions.HasFilePattern(),
			"hasLinearHistory":  pullRequest_functions.HasLinearHistory(),
			"hasLinkedIssues":   pullRequest_functions.HasLinkedIssues(),
			"head":              pullRequest_functions.Head(),
			"isDraft":           pullRequest_functions.IsDraft(),
			"labels":            pullRequest_functions.Labels(),
			"milestone":         pullRequest_functions.Milestone(),
			"reviewers":         pullRequest_functions.Reviewers(),
			"size":              pullRequest_functions.Size(),
			"title":             pullRequest_functions.Title(),
			// Organization
			"organization": organization_functions.Organization(),
			"team":         organization_functions.Team(),
			// User
			"totalCreatedPullRequests": user_functions.TotalCreatedPullRequests(),
			// Utilities
			"append":      utils_functions.AppendString(),
			"contains":    utils_functions.Contains(),
			"isElementOf": utils_functions.IsElementOf(),
			// Engine
			"group": engine_functions.Group(),
			"rule":  engine_functions.Rule(),
			// Internal
			"filter": inner_functions.Filter(),
		},
		Actions: map[string]*aladino.BuiltInAction{
			"addLabel":             actions.AddLabel(),
			"assignRandomReviewer": actions.AssignRandomReviewer(),
			"assignReviewer":       actions.AssignReviewer(),
			"assignTeamReviewer":   actions.AssignTeamReviewer(),
			"close":                actions.Close(),
			"comment":              actions.Comment(),
			"commentOnce":          actions.CommentOnce(),
			"merge":                actions.Merge(),
			"removeLabel":          actions.RemoveLabel(),
		},
	}
}
