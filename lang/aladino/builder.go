// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

func buildFilter(param string, condition Expr) (Expr, error) {
	organizationAST := functionCall(
		variable("organization"),
		[]Expr{},
	)

	ast := functionCall(
		variable("filter"),
		[]Expr{
			organizationAST,
			lambda(
				[]Expr{typedExpr(variable(param), BuildStringType())},
				condition,
			),
		},
	)

	return ast, nil
}
