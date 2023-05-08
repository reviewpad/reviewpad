// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
)

type ExecExpr interface {
	exec(env Env) error
}

func TypeCheckExec(env Env, expr Expr) (ExecExpr, error) {
	switch expr.Kind() {
	case "FunctionCall":
		_, err := TypeInference(env, expr)
		if err != nil {
			return nil, err
		}

		return expr.(*FunctionCall), nil
	}

	return nil, fmt.Errorf("typecheckexec: %v", expr.Kind())
}

func (fc *FunctionCall) exec(env Env) error {
	args := make([]Value, len(fc.arguments))
	for i, elem := range fc.arguments {
		value, err := elem.Eval(env)

		if err != nil {
			return err
		}

		args[i] = value
	}

	action, ok := env.GetBuiltIns().Actions[fc.name.ident]
	if !ok {
		return fmt.Errorf("exec: %v not found. are you sure this is a built-in function?", fc.name.ident)
	}

	if action.Disabled {
		env.GetLogger().Infof("action %v is disabled - skipping", fc.name.ident)
		return nil
	}

	collectedData := map[string]interface{}{
		"builtin": fc.name.ident,
	}

	if err := env.GetCollector().Collect("Ran Builtin", collectedData); err != nil {
		env.GetLogger().Errorf("error collection built-in run: %v\n", err)
	}

	entityKind := env.GetTarget().GetTargetEntity().Kind

	for _, supportedKind := range action.SupportedKinds {
		if entityKind == supportedKind {
			if action.RunAsynchronously {
				env.GetExecWaitGroup().Add(1)
				go func() {
					defer env.GetExecWaitGroup().Done()

					err := action.Code(env, args)
					if err != nil {
						env.SetExecFatalErrorOccurred(err)
					}
				}()
				return nil
			}
			return action.Code(env, args)
		}
	}

	return fmt.Errorf("eval: unsupported kind %v", entityKind)
}
