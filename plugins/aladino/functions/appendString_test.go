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

var appendString = plugins_aladino.PluginBuiltIns().Functions["append"].Code

func TestAppendString(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	slice := lang.BuildArrayValue(
		[]lang.Value{
			lang.BuildStringValue("a"),
			lang.BuildStringValue("b"),
		},
	)
	elems := lang.BuildArrayValue(
		[]lang.Value{
			lang.BuildStringValue("c"),
		},
	)
	args := []lang.Value{slice, elems}

	wantSlice := lang.BuildArrayValue(
		[]lang.Value{
			lang.BuildStringValue("a"),
			lang.BuildStringValue("b"),
			lang.BuildStringValue("c"),
		},
	)

	gotSlice, err := appendString(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantSlice, gotSlice)
}
