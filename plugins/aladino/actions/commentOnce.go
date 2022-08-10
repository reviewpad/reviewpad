// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/go-github/v45/github"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

const ReviewpadCommentAnnotation = "<!--@annotation-reviewpad-single-comment-->"

func CommentOnce() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: commentOnceCode,
	}
}

func commentOnceCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()

	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	commentBody := args[0].(*aladino.StringValue).Val
	commentBodyWithReviewpadAnnotation := fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, commentBody)
	commentBodyWithReviewpadAnnotationHash := sha256.Sum256([]byte(commentBodyWithReviewpadAnnotation))

	comments, err := e.GetGithubClient().GetPullRequestComments(e.GetCtx(), owner, repo, prNum, &github.IssueListCommentsOptions{})
	if err != nil {
		return err
	}

	for _, comment := range comments {
		commentHash := sha256.Sum256([]byte(*comment.Body))
		commentAlreadyExists := commentHash == commentBodyWithReviewpadAnnotationHash
		if commentAlreadyExists {
			return nil
		}
	}

	_, _, err = e.GetGithubClient().CreateComment(e.GetCtx(), owner, repo, prNum, &github.IssueComment{
		Body: &commentBodyWithReviewpadAnnotation,
	})

	return err
}
