// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/shurcooL/githubv4"
)

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
	clientGraphQL := e.GetGithubClient().GetClientGraphQL()

	reviewState, err := parseReviewState(args[0].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	reviewBody, err := parseReviewBody(reviewState, args[1].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	authenticatedUserLogin, err := getAuthenticatedUserLogin(clientGraphQL)
	if err != nil {
		return err
	}

	latestReview, err := getLatestReviewFromReviewer(clientGraphQL, entity, authenticatedUserLogin)
	if err != nil {
		return err
	}

	if latestReview != nil {
		log.Printf("review: latest review from %v is %v with body %v", authenticatedUserLogin, latestReview.State, latestReview.Body)
		if latestReview.State == reviewState && latestReview.Body == reviewBody {
			log.Printf("review: skipping review since it's the same as the latest review")
			return nil
		}
	}
	log.Printf("review: creating review %v with body %v", reviewState, reviewBody)

	return t.Review(reviewState, reviewBody)
}

func parseReviewState(reviewState string) (string, error) {
	switch reviewState {
	case "COMMENT", "REQUEST_CHANGES", "APPROVE":
		return reviewState, nil
	default:
		return "", fmt.Errorf("review: unsupported review state %v", reviewState)
	}
}

func parseReviewBody(reviewState, reviewBody string) (string, error) {
	if reviewState != "APPROVE" && reviewBody == "" {
		return "", fmt.Errorf("review: comment required in %v state", reviewState)
	}

	return reviewBody, nil
}

func getAuthenticatedUserLogin(clientGQL *githubv4.Client) (string, error) {
	var userLogin struct {
		Viewer struct {
			Login string
		}
	}

	err := clientGQL.Query(context.Background(), &userLogin, nil)
	if err != nil {
		return "", err
	}

	return userLogin.Viewer.Login, nil
}

func getLatestReviewFromReviewer(clientGQL *githubv4.Client, target *handler.TargetEntity, author string) (*codehost.Review, error) {
	var reviewsQuery struct {
		Repository struct {
			PullRequest struct {
				Reviews struct {
					Nodes []struct {
						Author struct {
							Login githubv4.String
						}
						Body        githubv4.String
						State       githubv4.String
						SubmittedAt *time.Time
					}
				} `graphql:"reviews(last: 1, author: $author)"`
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}
	varGQLReviews := map[string]interface{}{
		"repositoryOwner":   githubv4.String(target.Owner),
		"repositoryName":    githubv4.String(target.Repo),
		"pullRequestNumber": githubv4.Int(target.Number),
		"author":            author,
	}

	err := clientGQL.Query(context.Background(), &reviewsQuery, varGQLReviews)
	if err != nil {
		return nil, err
	}

	reviews := reviewsQuery.Repository.PullRequest.Reviews.Nodes
	if len(reviews) == 0 {
		return nil, nil
	}

	latestReview := reviews[0]
	return &codehost.Review{
		User: &codehost.User{
			Login: string(latestReview.Author.Login),
		},
		Body:        string(latestReview.Body),
		State:       string(latestReview.State),
		SubmittedAt: latestReview.SubmittedAt,
	}, nil
}
