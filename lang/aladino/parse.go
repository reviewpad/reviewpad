// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"strings"
)

func Parse(input string) (Expr, error) {
	input = strings.TrimRight(input, "\n")
	lex := &AladinoLex{input: input}
	res := AladinoParse(lex)

	if res != 0 {
		return nil, fmt.Errorf("parse error: failed to build AST on input %v", input)
	}

	return lex.ast, nil
}
