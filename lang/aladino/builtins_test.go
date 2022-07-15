// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
    "github.com/stretchr/testify/assert"
)

func TestMergeAladinoBuiltIns(t *testing.T) {
	builtInsList := []*aladino.BuiltIns{
		{
			Functions: map[string]*aladino.BuiltInFunction{
				"emptyFunction1": nil,
			},
			Actions: map[string]*aladino.BuiltInAction{
				"emptyAction1": nil,
			},
		},
		{
			Functions: map[string]*aladino.BuiltInFunction{
				"emptyFunction2": nil,
			},
			Actions: map[string]*aladino.BuiltInAction{
				"emptyAction2": nil,
			},
		},
	}

	wantBuiltIns := &aladino.BuiltIns{
		Functions: map[string]*aladino.BuiltInFunction{
			"emptyFunction1": nil,
			"emptyFunction2": nil,
		},
		Actions: map[string]*aladino.BuiltInAction{
			"emptyAction1": nil,
			"emptyAction2": nil,
		},
	}

	gotBuiltIns := aladino.MergeAladinoBuiltIns(builtInsList...)

	assert.Equal(t, wantBuiltIns, gotBuiltIns)
}
