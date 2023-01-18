// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

//go:build integration
// +build integration

package integration_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/reviewpad/reviewpad/v3"
	"github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	conflictCommitHead = "3e78768332b8f0c841bfaa9ed774d7c4b882c5b3"
)

type IntegrationTestMutation struct {
	CreateRef struct {
		ClientMutationID string
	} `graphql:"createRef(input: $input)"`
	CreateCommitOnBranch struct {
		ClientMutationID string
	} `graphql:"createCommitOnBranch(input: $createCommitAndBranchInput)"`
	CreatePullRequest struct {
		ClientMutationID string
		PullRequest      struct {
			ID     githubv4.ID
			Number int
		}
	} `graphql:"createPullRequest(input: $createPullRequestInput)"`
}

type UpdateIntegrationTestMutation struct {
	UpdatePullRequest struct {
		ClientMutationID string
	} `graphql:"updatePullRequest(input: $input)"`
}

type IDNode struct {
	ID githubv4.ID
}

type IntegrationTestQuery struct {
	Repository struct {
		ID               githubv4.ID
		DefaultBranchRef struct {
			Name string
		}
		Milestones struct {
			Nodes []IDNode
		} `graphql:"milestones(first: 1)"`
		Labels struct {
			Nodes []IDNode
		} `graphql:"labels(first: 3, query: $labelsQuery)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func TestIntegration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	ctx := context.Background()
	logger := logrus.NewEntry(logrus.StandardLogger())
	collector := engine.DefaultMockCollector
	githubToken := os.Getenv("GITHUB_INTEGRATION_TEST_TOKEN")
	repoOwner := os.Getenv("GITHUB_INTEGRATION_TEST_REPO_OWNER")
	repoName := os.Getenv("GITHUB_INTEGRATION_TEST_REPO_NAME")

	require.NotEqual(githubToken, "")
	require.NotEqual(repoOwner, "")
	require.NotEqual(repoName, "")

	githubClient := github.NewGithubClientFromToken(ctx, githubToken)

	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)
	branchName := fmt.Sprintf("integration-test-%s", uuid.NewString())
	addFile, err := utils.ReadFile("./assets/utils/add.go")
	require.Nil(err)
	subFile, err := utils.ReadFile("./assets/utils/sub.go")
	require.Nil(err)
	readmeFile, err := utils.ReadFile("./assets/README.md")
	require.Nil(err)
	binaryFile, err := utils.ReadFile("./assets/hello-world")
	require.Nil(err)
	rawMostBuiltInsReviewpadFile, err := utils.ReadFile("./assets/most-built-ins-reviewpad.yml")
	require.Nil(err)
	mostBuiltInsReviewpadFile, err := engine.Load(ctx, githubClient, rawMostBuiltInsReviewpadFile)
	require.Nil(err)

	var integrationTestQuery IntegrationTestQuery
	integrationTestQueryData := map[string]interface{}{
		"owner":       githubv4.String(repoOwner),
		"name":        githubv4.String(repoName),
		"labelsQuery": githubv4.String("bug | documentation | duplicate"),
	}

	err = githubClient.GetClientGraphQL().Query(ctx, &integrationTestQuery, integrationTestQueryData)
	require.Nil(err)
	require.NotEmpty(integrationTestQuery.Repository.ID)
	require.Len(integrationTestQuery.Repository.Milestones.Nodes, 1)
	require.Len(integrationTestQuery.Repository.Labels.Nodes, 3)

	tests := map[string]struct {
		createPullRequestInput githubv4.CreatePullRequestInput
		updatePullRequestInput githubv4.UpdatePullRequestInput
		reviewpadFile          *engine.ReviewpadFile
	}{
		"kitchen-sink": {
			createPullRequestInput: githubv4.CreatePullRequestInput{
				RepositoryID: integrationTestQuery.Repository.ID,
				BaseRefName:  githubv4.String(integrationTestQuery.Repository.DefaultBranchRef.Name),
				HeadRefName:  githubv4.String(fmt.Sprintf("refs/heads/%s", branchName)),
				Title:        githubv4.String(branchName),
				Body:         githubv4.NewString(githubv4.String("")),
				Draft:        githubv4.NewBoolean(githubv4.Boolean(true)),
			},
			updatePullRequestInput: githubv4.UpdatePullRequestInput{
				MilestoneID: &integrationTestQuery.Repository.Milestones.Nodes[0].ID,
				LabelIDs: mapSlice(integrationTestQuery.Repository.Labels.Nodes, func(node IDNode) githubv4.ID {
					return node.ID
				}),
			},
			reviewpadFile: mostBuiltInsReviewpadFile,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			createRefInput := githubv4.CreateRefInput{
				RepositoryID: githubv4.ID(integrationTestQuery.Repository.ID),
				Name:         githubv4.String(fmt.Sprintf("refs/heads/%s", branchName)),
				Oid:          githubv4.GitObjectID(conflictCommitHead),
			}
			onBoardMutationData := map[string]interface{}{
				"createCommitAndBranchInput": githubv4.CreateCommitOnBranchInput{
					Branch: githubv4.CommittableBranch{
						RepositoryNameWithOwner: githubv4.NewString(githubv4.String(repoFullName)),
						BranchName:              githubv4.NewString(githubv4.String(branchName)),
					},
					FileChanges: &githubv4.FileChanges{
						Additions: &[]githubv4.FileAddition{
							{
								Path:     githubv4.String("utils/add.go"),
								Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString(addFile)),
							},
							{
								Path:     githubv4.String("utils/sub.go"),
								Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString(subFile)),
							},
							{
								Path:     githubv4.String("hello-world"),
								Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString(binaryFile)),
							},
							{
								Path:     githubv4.String("README.md"),
								Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString(readmeFile)),
							},
						},
					},
					Message: githubv4.CommitMessage{
						Headline: githubv4.String("unconventional commit"),
					},
					ExpectedHeadOid: githubv4.GitObjectID(conflictCommitHead),
				},
				"createPullRequestInput": test.createPullRequestInput,
			}

			var integrationTestMutation IntegrationTestMutation
			err = githubClient.GetClientGraphQL().Mutate(ctx, &integrationTestMutation, createRefInput, onBoardMutationData)
			assert.Nil(err)

			test.updatePullRequestInput.PullRequestID = integrationTestMutation.CreatePullRequest.PullRequest.ID

			var updateIntegrationRequestMutation UpdateIntegrationTestMutation
			err = githubClient.GetClientGraphQL().Mutate(ctx, &updateIntegrationRequestMutation, test.updatePullRequestInput, nil)
			assert.Nil(err)

			targetEntity := &handler.TargetEntity{
				Kind:   handler.PullRequest,
				Owner:  repoOwner,
				Repo:   repoName,
				Number: integrationTestMutation.CreatePullRequest.PullRequest.Number,
			}

			eventDetails := &handler.EventDetails{
				EventName:   "pull_request",
				EventAction: "opened",
			}

			_, _, err := reviewpad.Run(ctx, logger, githubClient, collector, targetEntity, eventDetails, nil, test.reviewpadFile, false, false)
			assert.Nil(err)
		})
	}
}

func mapSlice[T any, M any](a []T, f func(T) M) *[]M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return &n
}
