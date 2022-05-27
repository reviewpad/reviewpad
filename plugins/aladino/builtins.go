// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import "github.com/reviewpad/reviewpad/lang/aladino"

func PluginBuiltIns() *aladino.BuiltIns {
	return &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			"addDevs":          allDevs(),
			"append":           appendString(),
			"assignees":        assignees(),
			"base":             base(),
			"codeQuery":        codeQuery(),
			"commits":          commits(),
			"commitsCount":     commitsCount(),
			"contains":         contains(),
			"createdAt":        createdAt(),
			"description":      description(),
			"fileCount":        fileCount(),
			"filesExtensions":  filesExtensions(),
			"filter":           filter(),
			"group":            group(),
			"hasFileName":      hasFileName(),
			"hasFilePattern":   hasFilePattern(),
			"hasLinearHistory": hasLinearHistory(),
			"hasLinkedIssues":  hasLinkedIssues(),
			"head":             head(),
			"isDraft":          isDraft(),
			"isMemberOf":       isMemberOf(),
			"labels":           labels(),
			"milestone":        milestone(),
			"name":             name(),
			"reviewers":        reviewers(),
			"size":             size(),
			"team":             team(),
			"title":            title(),
			"totalCreatedPRs":  totalCreatedPRs(),
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
