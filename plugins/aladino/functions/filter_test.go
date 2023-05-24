// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var filter = plugins_aladino.PluginBuiltIns().Functions["filter"].Code

func TestFilter(t *testing.T) {
	mockedIntValue := lang.BuildIntValue(1)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []lang.Value{
		lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("1"), mockedIntValue}),
		lang.BuildFunctionValue(func(args []lang.Value) lang.Value {
			return lang.BuildBoolValue(args[0].HasKindOf(lang.INT_VALUE))
		}),
	}
	gotElems, err := filter(mockedEnv, args)

	wantElems := lang.BuildArrayValue([]lang.Value{mockedIntValue})

	assert.Nil(t, err)
	assert.Equal(t, wantElems, gotElems)
}
