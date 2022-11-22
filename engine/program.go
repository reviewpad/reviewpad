// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import "github.com/reviewpad/cookbook/recipes"

type Statement struct {
	code string
}

func (s *Statement) GetStatementCode() string {
	return s.code
}

func BuildStatement(code string) *Statement {
	return &Statement{
		code,
	}
}

type Program struct {
	recipes    []recipes.Recipe
	statements []*Statement
}

func BuildProgram(recipes []recipes.Recipe, statements []*Statement) *Program {
	return &Program{
		recipes,
		statements,
	}
}

func (p *Program) GetStatements() []*Statement {
	return p.statements
}

func (p *Program) GetRecipes() []recipes.Recipe {
	return p.recipes
}

func (p *Program) appendInstructions(instructions []string) {
	for _, instruction := range instructions {
		statement := BuildStatement(instruction)
		p.statements = append(p.statements, statement)
	}
}
