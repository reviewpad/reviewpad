// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	actions "github.com/reviewpad/reviewpad/v2/plugins/aladino/actions"
	functions "github.com/reviewpad/reviewpad/v2/plugins/aladino/functions"
)

// The documentation for the builtins is in:
// https://github.com/reviewpad/docs/blob/main/aladino/builtins.md
// This means that changes to the builtins need to be propagated to that document.

func PluginBuiltIns() *aladino.BuiltIns {
	return &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			// Pull Request
			"assignees":         functions.Assignees(),
			"author":            functions.Author(),
			"base":              functions.Base(),
			"commitCount":       functions.CommitCount(),
			"commits":           functions.Commits(),
			"createdAt":         functions.CreatedAt(),
			"description":       functions.Description(),
			"fileCount":         functions.FileCount(),
			"hasCodePattern":    functions.HasCodePattern(),
			"hasFileExtensions": functions.HasFileExtensions(),
			"hasFileName":       functions.HasFileName(),
			"hasFilePattern":    functions.HasFilePattern(),
			"hasLinearHistory":  functions.HasLinearHistory(),
			"hasLinkedIssues":   functions.HasLinkedIssues(),
			"head":              functions.Head(),
			"isDraft":           functions.IsDraft(),
			"labels":            functions.Labels(),
			"milestone":         functions.Milestone(),
			"reviewers":         functions.Reviewers(),
			"size":              functions.Size(),
			"title":             functions.Title(),
			// Organization
			"organization": functions.Organization(),
			"team":         functions.Team(),
			// User
			"totalCreatedPullRequests": functions.TotalCreatedPullRequests(),
			// Utilities
			"append":      functions.AppendString(),
			"contains":    functions.Contains(),
			"isElementOf": functions.IsElementOf(),
			"startsWith":  functions.StartsWith(),
			// Engine
			"group": functions.Group(),
			"rule":  functions.Rule(),
			// Internal
			"filter": functions.Filter(),
		},
		Actions: map[string]*aladino.BuiltInAction{
			"addLabel":             actions.AddLabel(),
			"assignAssignees":      actions.AssignAssignees(),
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
