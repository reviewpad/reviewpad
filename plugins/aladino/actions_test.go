// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/lang/aladino"
	"github.com/stretchr/testify/assert"
)

const ReviewpadCommentAnnotation = "<!--@annotation-reviewpad-->"

func TestCommentOnceOnListCommentsFail(t *testing.T) {
	testEvalEnv, err := mockEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			&github.ErrorResponse{},
		),
	)
	if err != nil {
		log.Fatalf("mockEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue("Lorem Ipsum")}

	err = commentOnceCode(testEvalEnv, args)

	assert.NotNil(t, err)
}

func TestCommentOnceWhenCommentAlreadyExists(t *testing.T) {
	const ExistingComment = "Lorem Ipsum"
	const ExistingCommentWithReviewpadCommentAnnotation = ReviewpadCommentAnnotation + "Lorem Ipsum"
	
	var commentCreated *string
	testEvalEnv, err := mockEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String(ExistingCommentWithReviewpadCommentAnnotation),
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				postBody := github.IssueComment{}

				json.Unmarshal(body, &postBody)

				commentCreated = postBody.Body
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(ExistingComment)}

	err = commentOnceCode(testEvalEnv, args)
	assert.Nil(t, err)
	assert.Equal(t, (*string)(nil), commentCreated)
}

func TestCommentOnce(t *testing.T) {
	var commentCreated *string
	testEvalEnv, err := mockEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String(ReviewpadCommentAnnotation + "Lorem Ipsum"),
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				postBody := github.IssueComment{}

				json.Unmarshal(body, &postBody)

				commentCreated = postBody.Body
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockEnv failed: %v", err)
	}

	const NewComment = "Dummy Comment"
	const NewCommentWithReviewpadAnnotation = ReviewpadCommentAnnotation + NewComment
	args := []aladino.Value{aladino.BuildStringValue(NewComment)}

	err = commentOnceCode(testEvalEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, NewCommentWithReviewpadAnnotation, *commentCreated)
}
