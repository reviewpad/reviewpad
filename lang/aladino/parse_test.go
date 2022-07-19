// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse_WhenSingleLine(t *testing.T) {
	input := `$addLabel("small")`
	wantExpr := BuildFunctionCall(
		BuildVariable("addLabel"),
		[]Expr{BuildStringConst("small")},
	)

	gotExpr, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}

func TestParse_WhenMultilineLine(t *testing.T) {
	input := `$addLabel("medium multiline")

`
	wantExpr := BuildFunctionCall(
		BuildVariable("addLabel"),
		[]Expr{BuildStringConst("medium multiline")},
	)

	gotExpr, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}
