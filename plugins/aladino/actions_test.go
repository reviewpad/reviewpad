// Copyright (C) 2019-2022 Explore.dev - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

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

	commentOnce := PluginBuiltIns().Actions["commentOnce"].Code

	err = commentOnce(*testEvalEnv, args)

	assert.NotNil(t, err)
}

func TestCommentOnceWhenCommentAlreadyExists(t *testing.T) {
	var commentCreated string
	testEvalEnv, err := mockEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String("Lorem Ipsum"),
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber ,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				postBody := github.IssueComment{}

				json.Unmarshal(body, &postBody)

				commentCreated = *postBody.Body
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue("Lorem Ipsum")}

	commentOnce := PluginBuiltIns().Actions["commentOnce"].Code

	err = commentOnce(*testEvalEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, "", commentCreated)
}

func TestCommentOnce(t *testing.T) {
	var commentCreated string
	testEvalEnv, err := mockEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			[]*github.IssueComment{
				{
					Body: github.String("Lorem Ipsum"),
				},
			},
		),
		mock.WithRequestMatchHandler(
			mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber ,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				postBody := github.IssueComment{}

				json.Unmarshal(body, &postBody)

				commentCreated = *postBody.Body
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue("Dummy Comment")}

	commentOnce := PluginBuiltIns().Actions["commentOnce"].Code

	err = commentOnce(*testEvalEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, "Dummy Comment", commentCreated)
}