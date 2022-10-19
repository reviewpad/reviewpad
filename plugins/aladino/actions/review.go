// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"log"
	"reflect"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
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
	reviewer := getWorkflowRunActor(e).GetLogin()

	log.Printf("%+v", reviewer)

	reviewEvent, err := parseReviewEvent(args[0].(*aladino.StringValue).Val)
	if err != nil {
		return err
	}

	reviewBody, err := checkReviewBody(reviewEvent, args[1].(*aladino.StringValue).Val)

	reviews, err := t.GetReviews()
	if err != nil {
		return err
	}

	if codehost.HasReview(reviews, reviewer) {
		lastReview := codehost.LastReview(reviews, reviewer)

		if lastReview.State == "APPROVED" {
			return nil
		}

		log.Printf("%+v", t.PullRequest.GetUpdatedAt())
		if lastReview.SubmittedAt.After(t.PullRequest.GetUpdatedAt()) {
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

func getWorkflowRunActor(e aladino.Env) *github.User {
	log.Printf("%+v", e.GetGithubClient().GetClientREST().UserAgent)
	workflowPayload := e.GetEventPayload()
	if reflect.TypeOf(workflowPayload).String() != "*github.WorkflowRunEvent" {
		return nil
	}

	workflowRunPayload := workflowPayload.(*github.WorkflowRunEvent).WorkflowRun
	if workflowRunPayload == nil {
		return nil
	}

	return workflowRunPayload.Actor
}
