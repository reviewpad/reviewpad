// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var filter = plugins_aladino.PluginBuiltIns().Functions["filter"].Code

func TestFilter(t *testing.T) {
	mockedIntValue := aladino.BuildIntValue(1)
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("1"), mockedIntValue}),
		aladino.BuildFunctionValue(func(args []aladino.Value) aladino.Value {
			return aladino.BuildBoolValue(args[0].HasKindOf(aladino.INT_VALUE))
		}),
	}
	gotElems, err := filter(mockedEnv, args)

	wantElems := aladino.BuildArrayValue([]aladino.Value{mockedIntValue})

	assert.Nil(t, err)
	assert.Equal(t, wantElems, gotElems)
}
