// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/google/go-github/v52/github"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Merge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
		Code:           mergeCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func mergeCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	log := e.GetLogger().WithField("builtin", "merge")
	totalRetryCount := 2

	if t.PullRequest.Status != pbc.PullRequestStatus_OPEN || t.PullRequest.IsDraft {
		log.Infof("skipping action because pull request is not open or is a draft")
		return nil
	}

	mergeMethod, err := parseMergeMethod(args)
	if err != nil {
		return err
	}

	if len(e.GetBuiltInsReportedMessages()[aladino.SEVERITY_FAIL]) > 0 {
		return nil
	}

	// We need to check if the base branch of the pull request has GitHub Merge Queue enabled.
	// The queue is not enabled if we have nil entries.
	gitHubMergeQueueEntries, err := e.GetGithubClient().GetGitHubMergeQueueEntries(e.GetCtx(), t.PullRequest.GetBase().Repo.Owner, t.PullRequest.GetBase().Repo.Name, t.PullRequest.GetBase().Name, totalRetryCount)
	if err != nil {
		return err
	}

	isGitHubMergeQueueEnabled := gitHubMergeQueueEntries != nil

	if isGitHubMergeQueueEnabled {
		// We need to check if the pull request is already in the GitHub Merge Queue.
		// If it is, we don't need to do anything.
		for _, pullRequestNumber := range gitHubMergeQueueEntries {
			if int64(pullRequestNumber) == t.PullRequest.GetNumber() {
				log.Infof("skipping action because pull request is already in the GitHub Merge Queue")
				return nil
			}
		}
	}

	if e.GetCheckRunID() != nil {
		e.SetCheckRunConclusion("success")

		summary := "Reviewpad is about to merge this pull request"
		if isGitHubMergeQueueEnabled {
			summary = "Reviewpad is about to add this pull request to the GitHub Merge Queue"
		}

		if err := updateCheckRunWithSummary(e, summary); err != nil {
			return err
		}
	}

	if isGitHubMergeQueueEnabled {
		if err := e.GetGithubClient().AddPullRequestToGithubMergeQueue(e.GetCtx(), t.PullRequest.GetId()); err != nil {
			if e.GetCheckRunID() != nil {
				if err := updateCheckRunWithSummary(e, "The pull request cannot be added to the merge queue"); err != nil {
					return err
				}
			}

			e.GetLogger().WithError(err).Warnln("failed to add this pull request to the GitHub Merge Queue")
		}

		return nil
	}

	mergeErr := t.Merge(mergeMethod)
	if mergeErr != nil {
		if e.GetCheckRunID() != nil {
			err := updateCheckRunWithSummary(e, "The merge cannot be completed due to non-compliance with certain GitHub branch protection rules")
			if err != nil {
				return err
			}
		}

		e.GetLogger().WithError(mergeErr).Warnln("failed to merge pull request")
	}

	return nil
}

func updateCheckRunWithSummary(e aladino.Env, summary string) error {
	targetEntity := e.GetTarget().GetTargetEntity()
	_, _, err := e.GetGithubClient().GetClientREST().Checks.UpdateCheckRun(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, *e.GetCheckRunID(), github.UpdateCheckRunOptions{
		Name:       "reviewpad",
		Status:     github.String("completed"),
		Conclusion: github.String("success"),
		Output: &github.CheckRunOutput{
			Title:   github.String("Reviewpad about to merge"),
			Summary: github.String(summary),
		},
	})
	return err
}

func parseMergeMethod(args []lang.Value) (string, error) {
	if len(args) == 0 {
		return "merge", nil
	}

	mergeMethod := args[0].(*lang.StringValue).Val
	switch mergeMethod {
	case "merge", "rebase", "squash":
		return mergeMethod, nil
	default:
		return "", fmt.Errorf("merge: unsupported merge method %v", mergeMethod)
	}
}
