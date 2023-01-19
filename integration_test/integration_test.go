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
	// opening a pull request from this commit
	// will guarantee it to have a conflict
	conflictCommitHead = "f4973722e25d7e8bea265dffc496859c079ac2cb"
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

	testID := uuid.NewString()
	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)

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

	rawMergeReviewpadFile, err := utils.ReadFile("./assets/merge.yml")
	require.Nil(err)

	rawDeleteHeadBranchReviewpadFile, err := utils.ReadFile("./assets/delete-head-branch.yml")
	require.Nil(err)

	rawFailReviewpadFile, err := utils.ReadFile("./assets/fail.yml")
	require.Nil(err)

	rawCloseReviewpadFile, err := utils.ReadFile("./assets/close.yml")
	require.Nil(err)

	mostBuiltInsReviewpadFile, err := engine.Load(ctx, githubClient, rawMostBuiltInsReviewpadFile)
	require.Nil(err)

	mergeReviewpadFile, err := engine.Load(ctx, githubClient, rawMergeReviewpadFile)
	require.Nil(err)

	deleteHeadBranchReviewpadFile, err := engine.Load(ctx, githubClient, rawDeleteHeadBranchReviewpadFile)
	require.Nil(err)

	failReviewpadFile, err := engine.Load(ctx, githubClient, rawFailReviewpadFile)
	require.Nil(err)

	closeReviewpadFile, err := engine.Load(ctx, githubClient, rawCloseReviewpadFile)
	require.Nil(err)

	// contains a graphql query that fetches the necessary data
	// to run the integration tests
	var integrationTestQuery IntegrationTestQuery
	integrationTestQueryData := map[string]interface{}{
		"owner":       githubv4.String(repoOwner),
		"name":        githubv4.String(repoName),
		"labelsQuery": githubv4.String("bug | documentation | duplicate"),
	}

	err = githubClient.GetClientGraphQL().Query(ctx, &integrationTestQuery, integrationTestQueryData)
	// ensure the required things are present on the repo which are
	// at least one milestone
	// at least 3 pre configured labels named bug, documentation and duplicate.
	require.Nil(err)
	require.NotEmpty(integrationTestQuery.Repository.ID)
	require.Len(integrationTestQuery.Repository.Milestones.Nodes, 1)
	require.Len(integrationTestQuery.Repository.Labels.Nodes, 3)

	tests := map[string]struct {
		createPullRequestInput githubv4.CreatePullRequestInput
		updatePullRequestInput *githubv4.UpdatePullRequestInput
		reviewpadFiles         []*engine.ReviewpadFile
		fileChanges            *githubv4.FileChanges
		commitMessage          string
		wantErr                error
		exitStatus             []engine.ExitStatus
	}{
		"kitchen-sink": {
			createPullRequestInput: githubv4.CreatePullRequestInput{
				RepositoryID: integrationTestQuery.Repository.ID,
				BaseRefName:  githubv4.String(integrationTestQuery.Repository.DefaultBranchRef.Name),
				Body:         githubv4.NewString(githubv4.String("")),
				Draft:        githubv4.NewBoolean(githubv4.Boolean(true)),
			},
			updatePullRequestInput: &githubv4.UpdatePullRequestInput{
				MilestoneID: &integrationTestQuery.Repository.Milestones.Nodes[0].ID,
				LabelIDs: mapSlice(integrationTestQuery.Repository.Labels.Nodes, func(node IDNode) githubv4.ID {
					return node.ID
				}),
			},
			fileChanges: &githubv4.FileChanges{
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
			commitMessage:  "kitchen sink test",
			reviewpadFiles: []*engine.ReviewpadFile{mostBuiltInsReviewpadFile, deleteHeadBranchReviewpadFile},
			exitStatus:     []engine.ExitStatus{engine.ExitStatusSuccess, engine.ExitStatusSuccess},
		},
		"merge-and-delete": {
			createPullRequestInput: githubv4.CreatePullRequestInput{
				RepositoryID: integrationTestQuery.Repository.ID,
				BaseRefName:  githubv4.String(integrationTestQuery.Repository.DefaultBranchRef.Name),
				Body:         githubv4.NewString(githubv4.String("merge and delete head integration test")),
				Draft:        githubv4.NewBoolean(githubv4.Boolean(false)),
			},
			fileChanges: &githubv4.FileChanges{
				Additions: &[]githubv4.FileAddition{
					{
						Path:     githubv4.String(fmt.Sprintf("ids/%s", testID)),
						Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString([]byte(testID))),
					},
				},
			},
			commitMessage:  "test: merge and delete",
			reviewpadFiles: []*engine.ReviewpadFile{mergeReviewpadFile, deleteHeadBranchReviewpadFile},
			exitStatus:     []engine.ExitStatus{engine.ExitStatusSuccess, engine.ExitStatusSuccess},
		},
		"fail": {
			createPullRequestInput: githubv4.CreatePullRequestInput{
				RepositoryID: integrationTestQuery.Repository.ID,
				BaseRefName:  githubv4.String(integrationTestQuery.Repository.DefaultBranchRef.Name),
				Body:         githubv4.NewString(githubv4.String("fail")),
				Draft:        githubv4.NewBoolean(githubv4.Boolean(false)),
			},
			fileChanges: &githubv4.FileChanges{
				Additions: &[]githubv4.FileAddition{
					{
						Path:     githubv4.String("utils/add.go"),
						Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString(addFile)),
					},
					{
						Path:     githubv4.String("utils/sub.go"),
						Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString(subFile)),
					},
				},
			},
			commitMessage:  "test: fail",
			reviewpadFiles: []*engine.ReviewpadFile{failReviewpadFile, closeReviewpadFile, deleteHeadBranchReviewpadFile},
			exitStatus:     []engine.ExitStatus{engine.ExitStatusFailure, engine.ExitStatusSuccess, engine.ExitStatusSuccess},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testID := uuid.NewString()
			branchName := fmt.Sprintf("%s-integration-test-%s", name, testID)

			test.createPullRequestInput.HeadRefName = githubv4.String(branchName)
			test.createPullRequestInput.Title = githubv4.String(branchName)

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
					FileChanges: test.fileChanges,
					Message: githubv4.CommitMessage{
						Headline: githubv4.String(test.commitMessage),
					},
					ExpectedHeadOid: githubv4.GitObjectID(conflictCommitHead),
				},
				"createPullRequestInput": test.createPullRequestInput,
			}

			// create the pull request the reviewpad file will run against
			// and ensure it's successfully created
			var integrationTestMutation IntegrationTestMutation
			err = githubClient.GetClientGraphQL().Mutate(ctx, &integrationTestMutation, createRefInput, onBoardMutationData)
			require.Nil(err)

			// if any updates are required to the pull request after creation
			// perform the updates and ensure the update is performed successfully
			if test.updatePullRequestInput != nil {
				test.updatePullRequestInput.PullRequestID = integrationTestMutation.CreatePullRequest.PullRequest.ID

				var updateIntegrationRequestMutation UpdateIntegrationTestMutation
				err = githubClient.GetClientGraphQL().Mutate(ctx, &updateIntegrationRequestMutation, *test.updatePullRequestInput, nil)
				require.Nil(err)
			}

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

			// execute the reviewpad files one by one and
			// ensure there are no errors and exit statuses match
			for i, file := range test.reviewpadFiles {
				exitStatus, _, err := reviewpad.Run(ctx, logger, githubClient, collector, targetEntity, eventDetails, nil, file, false, false)
				assert.Equal(test.wantErr, err)
				assert.Equal(test.exitStatus[i], exitStatus)
			}
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
