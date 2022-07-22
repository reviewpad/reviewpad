// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"log"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

func TestProcessGroup_WhenGroupTypeFilterIsSetAndBuildGroupASTFails(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeFilter,
		"$group(\""+groupName+"\")",
		"dev",
		"$hasFileExtensions(",
	)

	assert.EqualError(t, err, "buildGroupAST: parse error: failed to build AST on input $hasFileExtensions(")
}

func TestProcessGroup_WhenGroupTypeFilterIsSet(t *testing.T) {
	member := "jane"
	ghMembers := []*github.User{
		{Login: github.String(member)},
	}
	mockedEnv, err := aladino.MockDefaultEnvWithBuiltIns(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetOrgsMembersByOrg,
				ghMembers,
			),
		},
		nil,
        plugins_aladino.PluginBuiltIns(),
	)
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeFilter,
		"$group(\""+groupName+"\")",
		"dev",
		"$hasFileExtensions([\".ts\"])",
	)

	gotVal := mockedEnv.GetRegisterMap()[groupName]

	wantVal := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue(member),
	})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSet(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		"[\"jane\"]",
		"",
		"",
	)

	gotVal := mockedEnv.GetRegisterMap()[groupName]

	wantVal := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("jane"),
	})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSetAndEvalGroupTypeInferenceFails(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		"$group(1)",
		"",
		"",
	)

	assert.EqualError(t, err, "type inference failed: mismatch in arg types on group")
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSetAndGroupExprIsNotAnArray(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		"1",
		"",
		"",
	)

	assert.EqualError(t, err, "expression is not a valid group")
}
