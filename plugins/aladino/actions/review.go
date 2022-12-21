// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"context"
	"fmt"
	"time"

	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/shurcooL/githubv4"
)

type GetAuthenticatedUserDataQuery struct {
	Viewer struct {
		Login string
	}
}

func Review() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, nil),
		Code:           reviewCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func reviewCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	entity := e.GetTarget().GetTargetEntity()

	var authenticatedUserData GetAuthenticatedUserDataQuery

	if err := e.GetGithubClient().GetClientGraphQL().Query(e.GetCtx(), &authenticatedUserData, nil); err != nil {
		return err
	}

	authenticatedUserLogin := authenticatedUserData.Viewer.Login

	reviewEvent, err := parseReviewEvent(args[0].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	reviewBody, err := checkReviewBody(reviewEvent, args[1].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	reviews, err := getReviews(e.GetGithubClient().GetClientGraphQL(), t)
	if err != nil {
		return err
	}

	if codehost.HasReview(reviews, authenticatedUserLogin) {
		lastReview := codehost.LastReview(reviews, authenticatedUserLogin)

		if lastReview.State == "APPROVED" {
			return nil
		}

		lastPushDate, err := e.GetGithubClient().GetPullRequestLastPushDate(e.GetCtx(), entity.Owner, entity.Repo, entity.Number)
		if err != nil {
			return err
		}

		if lastReview.SubmittedAt.After(lastPushDate) {
			return nil
		}
	}

	return t.Review(reviewEvent, reviewBody)
}

func parseReviewEvent(reviewEvent string) (string, error) {
	switch reviewEvent {
	case "COMMENT", "REQUEST_CHANGES", "APPROVE":
		return reviewEvent, nil
	default:
		return "", fmt.Errorf("review: unsupported review event %v", reviewEvent)
	}
}

func checkReviewBody(reviewEvent, reviewBody string) (string, error) {
	if reviewEvent != "APPROVE" && reviewBody == "" {
		return "", fmt.Errorf("review: comment required in %v event", reviewEvent)
	}

	return reviewBody, nil
}

type ReviewThreadsQuery struct {
	Repository struct {
		PullRequest struct {
			Reviews struct {
				Nodes []struct {
					Author struct {
						AvatarUrl    githubv4.URI
						Login        githubv4.String
						ResourcePath githubv4.URI
						Url          githubv4.URI
					}
					Body        githubv4.String
					State       githubv4.String
					SubmittedAt *time.Time
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"reviews(first: 10, after: $reviewsCursor)"`
		} `graphql:"pullRequest(number: $pullRequestNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

func getReviews(clientGQL *githubv4.Client, t *target.PullRequestTarget) ([]*codehost.Review, error) {
	var reviewThreadsQuery ReviewThreadsQuery
	reviews := make([]*codehost.Review, 0)
	hasNextPage := true
	retryCount := 2

	varGQLReviewThreads := map[string]interface{}{
		"repositoryOwner":   githubv4.String(t.GetTargetEntity().Owner),
		"repositoryName":    githubv4.String(t.GetTargetEntity().Repo),
		"pullRequestNumber": githubv4.Int(t.GetTargetEntity().Number),
		"reviewsCursor":     (*githubv4.String)(nil),
	}

	currentRequestRetry := 1

	for hasNextPage {
		err := clientGQL.Query(context.Background(), &reviewThreadsQuery, varGQLReviewThreads)
		if err != nil {
			currentRequestRetry++
			if currentRequestRetry <= retryCount {
				continue
			} else {
				return nil, err
			}
		} else {
			currentRequestRetry = 0
		}
		nodes := reviewThreadsQuery.Repository.PullRequest.Reviews.Nodes
		for _, node := range nodes {
			review := &codehost.Review{
				User: &codehost.User{
					Login: string(node.Author.Login),
				},
				Body:        string(node.Body),
				State:       string(node.State),
				SubmittedAt: node.SubmittedAt,
			}
			reviews = append(reviews, review)
		}
		hasNextPage = reviewThreadsQuery.Repository.PullRequest.Reviews.PageInfo.HasNextPage
		varGQLReviewThreads["reviewThreadsCursor"] = githubv4.NewString(reviewThreadsQuery.Repository.PullRequest.Reviews.PageInfo.EndCursor)
	}

	return reviews, nil
}
