// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var appendString = plugins_aladino.PluginBuiltIns().Functions["append"].Code

func TestAppendString(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	slice := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("a"),
			aladino.BuildStringValue("b"),
		},
	)
	elems := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("c"),
		},
	)
	args := []aladino.Value{slice, elems}

	wantSlice := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("a"),
			aladino.BuildStringValue("b"),
			aladino.BuildStringValue("c"),
		},
	)

	gotSlice, err := appendString(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantSlice, gotSlice)
}
