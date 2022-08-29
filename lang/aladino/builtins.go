// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import "github.com/reviewpad/reviewpad/v3/handler"

type BuiltIns struct {
	Functions map[string]*BuiltInFunction
	Actions   map[string]*BuiltInAction
	Services  map[string]interface{}
}

type BuiltInFunction struct {
	Type           Type
	Code           func(e Env, args []Value) (Value, error)
	SupportedKinds []handler.TargetEntityKind
}

type BuiltInAction struct {
	Type           Type
	Code           func(e Env, args []Value) error
	Disabled       bool
	SupportedKinds []handler.TargetEntityKind
}

func MergeAladinoBuiltIns(builtInsList ...*BuiltIns) *BuiltIns {
	mergedBuiltIns := &BuiltIns{
		Functions: map[string]*BuiltInFunction{},
		Actions:   map[string]*BuiltInAction{},
		Services:  map[string]interface{}{},
	}

	for _, builtIns := range builtInsList {
		for key, fn := range builtIns.Functions {
			mergedBuiltIns.Functions[key] = fn
		}

		for key, action := range builtIns.Actions {
			mergedBuiltIns.Actions[key] = action
		}

		for key, service := range builtIns.Services {
			mergedBuiltIns.Services[key] = service
		}
	}

	return mergedBuiltIns
}
