// Copyright (C) 2022 Tiago Ferreira - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import (
	"context"

	"github.com/google/go-github/v49/github"
	reviewpadGitHub "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/cookbook/ingredients"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Size struct {
	gitHubClient *reviewpadGitHub.GithubClient
	logger       *logrus.Entry
}

var sizeSmallLabel = ingredients.Label{
	Name:        "small",
	Color:       "219ebc",
	Description: "Pull request change is between 0 - 100 changes",
}

var sizeMediumLabel = ingredients.Label{
	Name:        "medium",
	Color:       "faedcd",
	Description: "Pull request change is between 101 - 500 changes",
}

var sizeLargeLabel = ingredients.Label{
	Name:        "large",
	Color:       "e76f51",
	Description: "Pull request change is more than 501+ changes",
}

type targetDetailsForSize struct {
	PrAdditions uint64
	PrDeletions uint64
	RepoLabels  []string
}

type targetDetailsQueryForSize struct {
	Repository struct {
		PullRequest struct {
			Additions uint64 `graphql:"additions"`
			Deletions uint64 `graphql:"deletions"`
		} `graphql:"pullRequest(number: $number)"`
		Labels struct {
			PageInfo struct {
				HasNextPage bool   `graphql:"hasNextPage"`
				EndCursor   string `graphql:"endCursor"`
			}
			Nodes []struct {
				Name string
			} `graphql:"nodes"`
		} `graphql:"labels(first: 100, after: $labelsEndCursor)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

func getTargetDetailsForSize(ctx context.Context, gitHubClient *reviewpadGitHub.GithubClient, targetEntity *handler.TargetEntity) (*targetDetailsForSize, error) {
	var targetDetails targetDetailsQueryForSize
	labels := make([]string, 0)
	hasNextPage := true
	variables := map[string]interface{}{
		"owner":           githubv4.String(targetEntity.Owner),
		"repo":            githubv4.String(targetEntity.Repo),
		"number":          githubv4.Int(targetEntity.Number),
		"labelsEndCursor": (*githubv4.String)(nil),
	}

	for hasNextPage {
		if err := gitHubClient.GetClientGraphQL().Query(ctx, &targetDetails, variables); err != nil {
			return nil, err
		}

		for _, node := range targetDetails.Repository.Labels.Nodes {
			labels = append(labels, node.Name)
		}

		variables["labelsEndCursor"] = targetDetails.Repository.Labels.PageInfo.EndCursor
		hasNextPage = targetDetails.Repository.Labels.PageInfo.HasNextPage
	}

	return &targetDetailsForSize{
		PrAdditions: targetDetails.Repository.PullRequest.Additions,
		PrDeletions: targetDetails.Repository.PullRequest.Deletions,
		RepoLabels:  labels,
	}, nil
}

func NewSizeRecipe(logger *logrus.Entry, gitHubClient *reviewpadGitHub.GithubClient) *Size {
	return &Size{
		gitHubClient: gitHubClient,
		logger:       logger,
	}
}

func (r *Size) Run(ctx context.Context, targetEntity handler.TargetEntity) error {
	targetDetails, err := getTargetDetailsForSize(ctx, r.gitHubClient, &targetEntity)
	if err != nil {
		return err
	}

	totalLinesChanged := targetDetails.PrAdditions + targetDetails.PrDeletions

	labelsToAdd := make([]ingredients.Label, 0)
	labelsToRemove := make([]ingredients.Label, 0)

	if totalLinesChanged <= 100 {
		labelsToAdd = append(labelsToAdd, sizeSmallLabel)
		labelsToRemove = append(labelsToRemove, sizeMediumLabel, sizeLargeLabel)
	} else if totalLinesChanged >= 100 && totalLinesChanged <= 500 {
		labelsToAdd = append(labelsToAdd, sizeMediumLabel)
		labelsToRemove = append(labelsToRemove, sizeSmallLabel, sizeLargeLabel)
	} else {
		labelsToAdd = append(labelsToAdd, sizeLargeLabel)
		labelsToRemove = append(labelsToRemove, sizeSmallLabel, sizeMediumLabel)
	}

	r.logger.Infof("adding labels: %v", labelsToAdd)
	r.logger.Infof("removing labels: %v", labelsToRemove)

	for _, label := range labelsToAdd {
		if !slices.Contains(targetDetails.RepoLabels, label.Name) {
			_, _, createLabelErr := r.gitHubClient.CreateLabel(ctx, targetEntity.Owner, targetEntity.Repo, &github.Label{
				Name:        &label.Name,
				Color:       &label.Color,
				Description: &label.Description,
			})
			if createLabelErr != nil {
				r.logger.WithError(err).Errorf("failed to create label: %s", label.Name)
			}
		}
	}

	labelToAddNames := make([]string, 0)
	for _, label := range labelsToAdd {
		labelToAddNames = append(labelToAddNames, label.Name)
	}
	_, _, err = r.gitHubClient.AddLabels(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number, labelToAddNames)
	if err != nil {
		r.logger.WithError(err).Errorf("failed to add labels: %v", labelToAddNames)
	}

	for _, label := range labelsToRemove {
		_, err = r.gitHubClient.RemoveLabelForIssue(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number, label.Name)
		if err != nil {
			r.logger.WithError(err).Errorf("failed to remove label: %s", label.Name)
		}
	}

	return nil
}
