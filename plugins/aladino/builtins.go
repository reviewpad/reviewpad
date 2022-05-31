// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import "github.com/reviewpad/reviewpad/lang/aladino"

func PluginBuiltIns() *aladino.BuiltIns {
	return &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			// Pull Request
			"assignees":              assignees(),
			"author":                 author(),
			"base":                   base(),
			"commitCount":            commitCount(),
			"commits":                commits(),
			"createdAt":              createdAt(),
			"description":            description(),
			"fileCount":              fileCount(),
			"hasLinearHistory":       hasLinearHistory(),
			"hasLinkedIssues":        hasLinkedIssues(),
			"head":                   head(),
			"isDraft":                isDraft(),
			"labels":                 labels(),
			"milestone":              milestone(),
			"patchHasCodePattern":    patchHasCodePattern(),
			"patchHasFileExtensions": patchHasFileExtensions(),
			"patchHasFileName":       patchHasFileName(),
			"patchHasFilePattern":    patchHasFilePattern(),
			"reviewers":              reviewers(),
			"size":                   size(),
			"title":                  title(),
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
			// Internal
			"filter": filter(),
		},
		Actions: map[string]*aladino.BuiltInAction{
			"addLabel":             addLabel(),
			"assignRandomReviewer": assignRandomReviewer(),
			"assignReviewer":       assignReviewer(),
			"merge":                merge(),
			"removeLabel":          removeLabel(),
		},
	}
}
