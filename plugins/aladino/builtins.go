// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"log"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	actions "github.com/reviewpad/reviewpad/v3/plugins/aladino/actions"
	functions "github.com/reviewpad/reviewpad/v3/plugins/aladino/functions"
	services "github.com/reviewpad/reviewpad/v3/plugins/aladino/services"
	"google.golang.org/grpc"
)

type PluginConfig struct {
	Services        map[string]interface{}
	grpcConnections []*grpc.ClientConn
}

func DefaultPluginConfig() (*PluginConfig, error) {
	connections := make([]*grpc.ClientConn, 0)
	semanticClient, semanticConnection, err := services.NewSemanticService()
	if err != nil {
		return nil, err
	}

	connections = append(connections, semanticConnection)

	services := map[string]interface{}{
		services.SEMANTIC_SERVICE_KEY: semanticClient,
	}

	config := &PluginConfig{
		Services:        services,
		grpcConnections: connections,
	}

	return config, nil
}

func (config *PluginConfig) CleanupPluginConfig() {
	for _, connection := range config.grpcConnections {
		connection.Close()
	}
}

// The documentation for the builtins is in:
// https://github.com/reviewpad/docs/blob/main/aladino/builtins.md
// This means that changes to the builtins need to be propagated to that document.
func PluginBuiltInsWithConfig(config *PluginConfig) *aladino.BuiltIns {
	return &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			// Pull Request
			"assignees":             functions.Assignees(),
			"author":                functions.Author(),
			"base":                  functions.Base(),
			"changed":               functions.Changed(),
			"commentCount":          functions.CommentCount(),
			"comments":              functions.Comments(),
			"commitCount":           functions.CommitCount(),
			"commits":               functions.Commits(),
			"createdAt":             functions.CreatedAt(),
			"description":           functions.Description(),
			"fileCount":             functions.FileCount(),
			"hasAnnotation":         functions.HasAnnotation(),
			"hasCodePattern":        functions.HasCodePattern(),
			"hasFileExtensions":     functions.HasFileExtensions(),
			"hasFileName":           functions.HasFileName(),
			"hasFilePattern":        functions.HasFilePattern(),
			"hasGitConflicts":       functions.HasGitConflicts(),
			"hasLinearHistory":      functions.HasLinearHistory(),
			"hasLinkedIssues":       functions.HasLinkedIssues(),
			"hasUnaddressedThreads": functions.HasUnaddressedThreads(),
			"head":                  functions.Head(),
			"isDraft":               functions.IsDraft(),
			"isWaitingForReview":    functions.IsWaitingForReview(),
			"labels":                functions.Labels(),
			"lastEventAt":           functions.LastEventAt(),
			"milestone":             functions.Milestone(),
			"requestedReviewers":    functions.RequestedReviewers(),
			"reviewers":             functions.Reviewers(),
			"reviewerStatus":        functions.ReviewerStatus(),
			"size":                  functions.Size(),
			"title":                 functions.Title(),
			"workflowStatus":        functions.WorkflowStatus(),
			// Organization
			"organization": functions.Organization(),
			"team":         functions.Team(),
			// User
			"issueCountBy":             functions.IssueCountBy(),
			"pullRequestCountBy":       functions.PullRequestCountBy(),
			"totalCodeReviews":         functions.TotalCodeReviews(),
			"totalCreatedPullRequests": functions.TotalCreatedPullRequests(),
			// Utilities
			"all":         functions.All(),
			"any":         functions.Any(),
			"append":      functions.AppendString(),
			"contains":    functions.Contains(),
			"isElementOf": functions.IsElementOf(),
			"length":      functions.Length(),
			"sprintf":     functions.Sprintf(),
			"startsWith":  functions.StartsWith(),
			// Engine
			"group": functions.Group(),
			"rule":  functions.Rule(),
			// Internal
			"filter": functions.Filter(),
		},
		Actions: map[string]*aladino.BuiltInAction{
			"addLabel":             actions.AddLabel(),
			"addToProject":         actions.AddToProject(),
			"assignAssignees":      actions.AssignAssignees(),
			"assignRandomReviewer": actions.AssignRandomReviewer(),
			"assignReviewer":       actions.AssignReviewer(),
			"assignTeamReviewer":   actions.AssignTeamReviewer(),
			"close":                actions.Close(),
			"comment":              actions.Comment(),
			"commentOnce":          actions.CommentOnce(),
			"commitLint":           actions.CommitLint(),
			"deleteHeadBranch":     actions.DeleteHeadBranch(),
			"disableActions":       actions.DisableActions(),
			"error":                actions.ErrorMsg(),
			"fail":                 actions.Fail(),
			"info":                 actions.Info(),
			"merge":                actions.Merge(),
			"rebase":               actions.Rebase(),
			"removeLabel":          actions.RemoveLabel(),
			"removeLabels":         actions.RemoveLabels(),
			"review":               actions.Review(),
			"titleLint":            actions.TitleLint(),
			"warn":                 actions.Warn(),
		},
		Services: config.Services,
	}
}

func PluginBuiltIns() *aladino.BuiltIns {
	config, err := DefaultPluginConfig()
	if err != nil {
		log.Fatal("Error loading default plugin config: ", err)
	}
	return PluginBuiltInsWithConfig(config)
}
