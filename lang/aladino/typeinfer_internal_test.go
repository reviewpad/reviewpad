// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeInference_WhenGivenBoolConst(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

    expr := BuildBoolConst(true)

    wantType := BuildBoolType()

    gotType, err := TypeInference(mockedEnv, expr)

    assert.Nil(t, err)
    assert.Equal(t, wantType, gotType, "bool type is expected")
}

func TestTypeInference_ExpectTypeInferenceFailed(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

    expr := BuildVariable("$nonBuiltIn")

    gotType, err := TypeInference(mockedEnv, expr)

    assert.Nil(t, gotType)
    assert.EqualError(t, err, "no type for built-in $nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it.")
}
