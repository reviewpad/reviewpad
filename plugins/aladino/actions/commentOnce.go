// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"crypto/sha256"
	"fmt"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

const ReviewpadCommentAnnotation = "<!--@annotation-reviewpad-single-comment-->"

func CommentOnce() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           commentOnceCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func commentOnceCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget()

	commentBody := args[0].(*lang.StringValue).Val
	commentBodyWithReviewpadAnnotation := fmt.Sprintf("%s\n%s", ReviewpadCommentAnnotation, commentBody)
	commentBodyWithReviewpadAnnotationHash := sha256.Sum256([]byte(commentBodyWithReviewpadAnnotation))

	comments, err := t.GetComments()
	if err != nil {
		return err
	}

	for _, comment := range comments {
		commentHash := sha256.Sum256([]byte(comment.Body))
		commentAlreadyExists := commentHash == commentBodyWithReviewpadAnnotationHash
		if commentAlreadyExists {
			return nil
		}
	}

	return t.Comment(commentBodyWithReviewpadAnnotation)
}
