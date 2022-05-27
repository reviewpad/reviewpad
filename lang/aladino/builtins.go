// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

type BuiltIns struct {
	Functions map[string]*BuiltInFunction
	Actions   map[string]*BuiltInAction
}

type BuiltInFunction struct {
	Type Type
	Code func(e *EvalEnv, args []Value) (Value, error)
}

type BuiltInAction struct {
	Type Type
	Code func(e *EvalEnv, args []Value) error
}
