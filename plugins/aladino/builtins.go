// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import "github.com/reviewpad/reviewpad/v2/lang/aladino"

// The documentation for the builtins is in:
// https://github.com/reviewpad/docs/blob/main/aladino/builtins.md
// This means that changes to the builtins need to be propagated to that document.

func PluginBuiltIns() *aladino.BuiltIns {
	return &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			// Pull Request
			"assignees":         assignees(),
			"author":            author(),
			"base":              base(),
			"commitCount":       commitCount(),
			"commits":           commits(),
			"createdAt":         createdAt(),
			"description":       description(),
			"fileCount":         fileCount(),
			"hasCodePattern":    hasCodePattern(),
			"hasFileExtensions": hasFileExtensions(),
			"hasFileName":       hasFileName(),
			"hasFilePattern":    hasFilePattern(),
			"hasLinearHistory":  hasLinearHistory(),
			"hasLinkedIssues":   hasLinkedIssues(),
			"head":              head(),
			"isDraft":           isDraft(),
			"labels":            labels(),
			"milestone":         milestone(),
			"reviewers":         reviewers(),
			"size":              size(),
			"title":             title(),
			// Organization
			"organization": organization(),
			"team":         team(),
			// User
			"totalCreatedPullRequests": totalCreatedPullRequests(),
			// Utilities
			"append":      appendString(),
			"contains":    contains(),
			"isElementOf": isElementOf(),
			// Engine
			"group": group(),
			"rule":  rule(),
			// Internal
			"filter": filter(),
		},
		Actions: map[string]*aladino.BuiltInAction{
			"addLabel":             addLabel(),
			"assignRandomReviewer": assignRandomReviewer(),
			"assignReviewer":       assignReviewer(),
			"close":                close(),
			"comment":              comment(),
			"commentOnce":          commentOnce(),
			"merge":                merge(),
			"removeLabel":          removeLabel(),
		},
	}
}
