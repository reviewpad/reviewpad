// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import "github.com/reviewpad/reviewpad/v4/lang"

func BuildFilter(param string, condition Expr) (Expr, error) {
	organizationAST := BuildFunctionCall(
		BuildVariable("organization"),
		[]Expr{},
	)

	ast := BuildFunctionCall(
		BuildVariable("filter"),
		[]Expr{
			organizationAST,
			BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable(param), lang.BuildStringType())},
				condition,
			),
		},
	)

	return ast, nil
}
