// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDefaultTotalRequestedReviewers_WhenRegexMatchNotFound(t *testing.T) {
	wantStr := "$addLabel(\"dummy\")"

	gotStr := addDefaultTotalRequestedReviewers(wantStr)

	assert.Equal(t, wantStr, gotStr)
}

func TestAddDefaultTotalRequestedReviewers_WhenRegexMatchFound(t *testing.T) {
	str := "$assignReviewer([\"john\", \"jane\"])"
	wantStr := "$assignReviewer([\"john\", \"jane\"], 99)"

	gotStr := addDefaultTotalRequestedReviewers(str)

	assert.Equal(t, wantStr, gotStr)
}

func TestAddDefaultMergeMethod_WhenStringDoesNotContainMergeAction(t *testing.T) {
	wantStr := "$addLabel(\"dummy\")"

	gotStr := addDefaultMergeMethod(wantStr)

	assert.Equal(t, wantStr, gotStr)
}

func TestAddDefaultMergeMethod_WhenStringContainsMergeAction(t *testing.T) {
	str := "$merge()"
	wantStr := "$merge(\"merge\")"

	gotStr := addDefaultMergeMethod(str)

	assert.Equal(t, wantStr, gotStr)
}

func TestTransformActionStr(t *testing.T) {
	str := "$assignReviewer([\"john\", \"jane\"]) $merge()"
	wantStr := "$assignReviewer([\"john\", \"jane\"], 99) $merge(\"merge\")"

	gotStr := transformActionStr(str)

	assert.Equal(t, wantStr, gotStr)
}
