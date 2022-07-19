// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
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

	prNum := utils.GetPullRequestNumber(pullRequest)
	owner := utils.GetPullRequestBaseOwnerName(pullRequest)
	repo := utils.GetPullRequestRepoName(pullRequest)

	commentBody := args[0].(*aladino.StringValue).Val
	commentBodyWithReviewpadAnnotation := fmt.Sprintf("%v%v", ReviewpadCommentAnnotation, commentBody)
	commentBodyWithReviewpadAnnotationHash := sha256.Sum256([]byte(commentBodyWithReviewpadAnnotation))

	comments, err := utils.GetPullRequestComments(e.GetCtx(), e.GetClient(), owner, repo, prNum, &github.IssueListCommentsOptions{})
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

	_, _, err = e.GetClient().Issues.CreateComment(e.GetCtx(), owner, repo, prNum, &github.IssueComment{
		Body: &commentBodyWithReviewpadAnnotation,
	})

	return err
}
